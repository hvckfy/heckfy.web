local space = box.schema.space.create('saveshare', {
    if_not_exists = true
})

-- Define field indices
local SAVESHARE_HASH = 1
local SAVESHARE_DATA = 2

-- Define the format for the space
space:format({
    { name = 'hash', type = 'string' },
    { name = 'data', type = 'string' }
})

-- Create a primary index on the 'hash' field
space:create_index('hashindex', {
    parts = { 'hash' },
    type = 'hash',
    if_not_exists = true
})

-- Function to delete used data by hash
function delete_used_data(hash)
    local tuple = space:get(hash) -- Get the tuple using the hash directly
    if tuple then
        space:delete(tuple[1]) -- Delete the tuple from space
        return "200" -- Success
    end
    return "404 Not Found" -- Not found
end

-- Function to add new data to the space
function add_new_data(data, hash)
    local data_tuple = {hash, data}
    
    -- Check if the record with the given hash already exists
    if space:get(hash) ~= nil then 
        return "200" -- Record already exists, but we consider this a successful operation
    end
    -- Attempt to insert the new tuple
    space:insert(data_tuple) 
    return "200" -- Success
end

-- Function to get data by hash
function get_data(hash)
    local tuple = space:get(hash) -- Get the tuple using the hash directly
    if tuple then
        return "200", tuple -- Return success and the tuple
    end
    return "404 Not Found" -- Not found
end
