package models

const IpMatchModel = `
[request_definition]
r = role, obj, act

[policy_definition]
p = role, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = ipMatch(r.role, p.role) && r.obj == p.obj && r.act == p.act
`
