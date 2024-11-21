-- Configuration for Tarantool
box.cfg{
    listen = 3307,
}

-- Start the console on localhost
require('console').listen('127.0.0.1:3317')

-- Create user and grant permissions
box.schema.user.create('heckfy', {
    password = 'pass',
    if_not_exists = true
})
box.schema.user.grant('heckfy', 'read,write,execute,create,drop', 'universe', nil, {
    if_not_exists = true
})

local fiber = require('fiber')

local space = box.schema.space.create('requests', {
    if_not_exists = true
})

local RQL_USER = 1
local RQL_TIMESTAMP = 2
local RQL_REQUESTS = 3

space:format({{

    name = 'user',
    type = 'string'
}, {
    name = 'timestamp',
    type = 'unsigned'
}, {
    name = 'requests',
    type = 'unsigned'
}})

-- primary idindex
space:create_index('userindex', {
    parts = {'user'},
    type = 'hash',
    if_not_exists = true
})

local function time_allow(user_time)
    return (math.floor(fiber.time()) - user_time) > 60
end

function access(user_ip)
    local user_tuple = box.space.requests.index.userindex:get(user_ip)

    if not user_tuple then
        -- User does not exist, create a new entry
        user_tuple = {
            [RQL_USER] = user_ip,
            [RQL_TIMESTAMP] = math.floor(fiber.time()), -- add here 1 min
            [RQL_REQUESTS] = 1
        }
        local status, err = box.space.requests:insert(user_tuple)
        if not status then
            return "500", false -- Handle insertion error
        end
        return "201", true -- Created a new user
    end

    -- Check if the time limit has expired
    if time_allow(user_tuple[RQL_TIMESTAMP]) then
        -- Reset the request count if 1 minute has expired
        local status, err = box.space.requests:update(user_tuple[RQL_USER], {
            {'=', RQL_TIMESTAMP, math.floor(fiber.time())},
            {'=', RQL_REQUESTS, 1}  -- Reset requests to 1
        })
        if not status then
            return "500", false -- Handle update error
        end
        return "202", true -- User updated
    else
        -- Handle rate limiting
        if user_tuple[RQL_REQUESTS] < 5 then
            -- Increment requests count
            local status, err = box.space.requests:update(user_tuple[RQL_USER], {
                {'+', RQL_REQUESTS, 1}  -- Increment request count
            })
            if not status then
                return "500", false -- Handle update error
            end
            return "200", true -- OK
        else
            -- Access denied
            return "429", false -- Too Many Requests
        end
    end
end
