package models

const DefaultModel = `
[request_definition]
r = role, obj, act

[policy_definition]
p = role, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.role == p.role && keyMatch5(r.obj, p.obj) && regexMatch(r.act, p.act)
`
