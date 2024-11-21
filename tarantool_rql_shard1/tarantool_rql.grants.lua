box.schema.user.create('LOGIN', {
    password = 'PASSWORD',
    if_not_exists = true
})
box.schema.user.grant('LOGIN', 'read,write,execute,create,drop', 'universe', nil, {
    if_not_exists = true
})