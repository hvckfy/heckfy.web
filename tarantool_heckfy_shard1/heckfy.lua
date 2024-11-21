local fiber = require('fiber')
local uuid = require('uuid')
local argon2 = require('argon2')

local space = box.schema.space.create('links', {
    if_not_exists = true
})

local LINKS_ID = 1
local LINKS_USER = 2
local LINKS_HASH = 3
local LINKS_LINK = 4
local LINKS_TIMESTAMP = 5

space:format({{
    name = 'id',
    type = 'string'
}, {
    name = 'user',
    type = 'string'
}, {
    name = 'hash',
    type = 'string'
}, {
    name = 'link',
    type = 'string'
}, {
    name = 'timestamp',
    type = 'unsigned'
}})

-- primary idindex
space:create_index('idindex', {
    parts = {'id'},
    type = 'hash',
    if_not_exists = true
})
-- hash hashindex
space:create_index('hashindex', {
    parts = {'hash'},
    type = 'hash',
    if_not_exists = true
})

-- Get_hash("123","123","123")
function Get_hash(userid, hash, link)

    local link_tuple, statuscode = Check_if_hash_exists(hash)

    if link_tuple ~= nil then
        return hash, statuscode
    else
        Add_new_link(userid, hash, link)
        return hash, "201" -- return code 201 after creating a new line
    end

end

function Get_link(userid, hash)
    -- check here if we have hash like this:

    local link_tuple, statuscode = Check_if_hash_exists(hash)

    if link_tuple ~= nil then
        return link_tuple[LINKS_LINK], statuscode -- return link with statuscode
    else
        return hash, "204" -- return code 204 - "no content", no link with this url found in db
    end
end

-- Add_new_link("123","123","123")
function Add_new_link(userid, hash, link)

    local curr_stamp = math.floor(fiber.time())

    local new_link = {
        [LINKS_ID] = uuid.str(),
        [LINKS_USER] = userid,
        [LINKS_HASH] = hash,
        [LINKS_LINK] = link,
        [LINKS_TIMESTAMP] = curr_stamp
    }

    return box.space.links:insert(new_link), err
end

function Check_if_hash_exists(hash)
    local link_tuple = box.space.links.index.hashindex:get(hash)
    if link_tuple ~= nil then
        return link_tuple, "200"
    else
        return nil, "204"
    end

    -- yes:return hash,"200" - ok

    -- no: return nil,"204" - no content

end
