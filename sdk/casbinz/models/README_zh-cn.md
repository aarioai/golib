# model 说明

matchers 可以用到的[函数](https://casbin.org/docs/function)

> bool keyMatch...|globMatch(url:string, pattern:string)
> string keyGet|keyGet2|keyGet3(url:string, pattern:string, key_name:string)

``` 
[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && r.act == p.act
```

| 函数         | r.obj                     | p.obj                      | key_name | result    |
|------------|---------------------------|----------------------------|----------|-----------|
| keyMatch   | /alice_data/resource1     | /alice_data/*              |          | true      |
| keyMatch2  | /alice_data/resource1     | /alice_data/:resource      |          | true      |
| keyMatch3  | /alice_data/resource1     | /alice_data/{resource}     |          | true      |
| keyMatch4  | /alice_data/123/book/123  | /alice_data/{id}/book/{id} |          | true      |
| keyMatch5  | /alice_data/123/?status=1 | /alice_data/{id}/*         |          | true      |
| regexMatch | r.act = GET               | p.act = (GET)\|(POST)      |          | true      |
| regexMatch | /topic/edit/123           | /topic/edit/\[0-9\]+       |          | true      |
| ipMatch    | 192.168.2.123             | 192.168.2.0/24             |          | true      |
| globMatch  | /alice_data/resource1     | /alice_data/*              |          | true      |
| keyGet     | /proj/resource1           | /proj/*                    |          | resource1 |
| keyGet2    | /proj/resource1           | /proj/:res                 | res      | resource1 |
| keyGet3    | /proj/res3_admin/         | /proj/{resource}_admin/*   | resource | res3      |

参考  [match operators test](https://github.com/casbin/casbin/blob/master/util/builtin_operators_test.go)