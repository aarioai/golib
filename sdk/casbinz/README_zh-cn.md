# Casbin 管理人员访问控制
Casbin：这是一个用Go语言编写的开源访问控制库，支持多种访问控制模型，包括RBAC（基于角色的访问控制）。它不仅支持操作权限，还支持数据权限，并且可以通过配置实现权限继承。

[github](https://github.com/casbin/casbin)




| 技术	             | 主要功能	          | 优点	                       | 缺点                       | 	适用场景             |
|-----------------|----------------|---------------------------|--------------------------|-------------------|
| **RBAC**        | 基于角色的访问控制      | 简单易用，权限管理集中，角色可重用         | 对于复杂权限需求灵活性不足，角色数量增加管理复杂 | 中小型企业，权限需求相对简单的系统 | 
| **RBAC & IAM**  | 身份和访问管理，包括RBAC | 综合管理用户身份和访问权限，提供全面的安全管理功能 | 复杂性较高，需要专业的安全管理人员	       | 大型企业，需要全面安全管理的系统  |
| **ABAC**	       | 基于属性的访问控制	     | 高度灵活，支持复杂访问控制需求，动态调整权限    | 配置复杂，学习曲线陡峭，性能开销较大       | 大型企业，需要精细权限控制的系统  | 
| ACL	            | 访问控制列表         | 	简单直观，直接控制资源访问	           | 管理复杂性上升，不支持复杂权限继承和角色管理	  | 小型系统，权限需求简单的系统    |
| PBAC            | 基于策略的访问控制      | 高度灵活，支持复杂访问控制逻辑，动态调整权限	   | 配置复杂，学习曲线陡峭，性能开销较大       | 大型企业，需要精细权限控制的系统  | 


casbin 由两部分组成： model.conf（或其他配置方式） 和 policy.csv（或其他数据库 adapter）

[模型语法](https://casbin.org/zh/docs/syntax-for-models/)

**request_definition** 定义了e.Enforce(...) 函数中的参数。
```conf
[request_definition]
r = sub, obj, act
```
在这个例子中，sub，obj和act代表了经典的访问三元组：主体（访问实体），对象（被访问资源）和动作（访问方法）。 然而，你可以自定义你自己的请求格式。 例如，如果你不需要指定特定的资源，你可以使用sub, act，或者如果你有两个访问实体，你可以使用sub, sub2, obj, act。

**policy_definition** 定义了策略的含义。 例如，我们有以下模型：
```conf
[policy_definition]
p = sub, obj, act
p2 = sub, act
```
我们有以下策略（如果在策略文件中）：
p, alice, data1, read
p2, bob, write-all-objects

**policy_effect** 策略效果的定义。 如果多个策略规则匹配请求，它决定是否应批准访问请求。 例如，一条规则允许，另一条规则拒绝。
```conf
[policy_effect]
e = some(where (p.eft == allow))
```
上述策略效果意味着，如果有任何匹配的allow策略规则，最终效果是allow（也称为允许覆盖）。 p.eft是策略的效果，它可以是allow或deny。 这是可选的，其默认值是allow。 由于我们在上面没有指定它，所以它使用默认值。

策略效果的另一个例子是：
```conf
[policy_effect]
e = !some(where (p.eft == deny))
```
这意味着，如果没有匹配的deny策略规则，最终的效果是allow（也称为deny-override）。 some意味着存在一个匹配的策略规则。 any意味着所有匹配的策略规则（这里未使用）。 策略效果甚至可以与逻辑表达式连接：

```conf
[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
```
这意味着必须至少有一个匹配的allow策略规则，并且不能有任何匹配的deny策略规则。 因此，以这种方式，支持允许和拒绝授权，并且拒绝优先。

支持的内置策略效果有：

| 策略效果                                                         | 含义               | 	    例子         |  
|--------------------------------------------------------------|------------------|-----------------| 
| some(where (p.eft == allow))                                 | 	allow-override	 | ACL, RBAC, etc. |                       
| !some(where (p.eft == deny))                                 | 	deny-override	  | Deny-override   |                            
| some(where (p.eft == allow)) && !some(where (p.eft == deny)) | 	allow-and-deny  | 	Allow-and-deny |
| priority(p.eft) \|\| deny	                                   | priority	        | Priority        |
| subjectPriority(p.eft)	                                      | 基于角色的优先级         | 	主题-优先级         |

matchers 是策略匹配器的定义。 匹配器是定义如何根据请求评估策略规则的表达式。
```conf
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```
上述匹配器是最简单的，意味着请求中的主题、对象和动作应与策略规则中的相匹配。 匹配器中可以使用算术运算符如+, -, *, /和逻辑运算符如&&, ||, !。

注意：应该将更耗时的表达式或函数放到后面判断，这样能大幅度提升执行效率

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```

```policy.csv
p, alice, data1, read
p, bob, data2, write
```


## 存储

model 可以使用 model.conf，通过 `casbin.NewEnforcer("model.conf", ...)` 读取，也可以使用字符串，通过 model.NewModelFromString(text) 读取。

policy 可以通过csv 存储，也可以使用 adapter API。通用的 [adapters](https://casbin.org/docs/adapters)

adapter interface:
* LoadPolicy() 从存储器（如MySQL)中读取policy rules
* SavePolicy() 保存
* AddPolicy() 增加
* RemovePolicy() 删除
* RemoveFilteredPolicy 移除匹配的

数据库结构：
* id 主键
* ptype: It corresponds to p, g, g2, etc.
* v0-v5: policy.csv 的6列（每列都是字符串）  --> 在用户表将uid赋予

| id	 | ptype | 	v0           | 	v1	        | v2      | 	v3	 | v4	 | v5 |
|-----|-------|---------------|-------------|---------|------|-----|----|
| 1   | 	p    | 	data2_admin	 | data2       | read	   |
| 2	  | p     | data2_admin   | data2	      | write		 |
| 3   | g	    | alice         | data2_admin |         | 


### 扩展存储

* casbin_role 角色（群组）：id:uint, pid:uint, v0:string, name:string, effective_at, expire_at, created_at, updated_at
* casbin_object 对象（即v1-v4）：id:uint, v:string（即v1-v4值）, name:string, created_at, updated_at
* casbin_policy：id:uint, ptype, v0, v1, v2, v3, v4, v5
* casbin_user: uid:uint64, roles:[]uint, created_at, updated_at

由上面生成 policy.csv官方模板来


#### casbin_role
一个群，只能有1个父群。子群成员也会被同时赋予父群成员。

| id | pid | v0     | 	name	 | 	ptype | effective_at	 | expire_at	 | created_at | updated_at |
|----|-----|--------|--------|--------|---------------|------------|------------|------------| 
| 1  | 0   | system | 系统管理员  | g      |
| 2  | 0   | admin  | 管理群	   | g      |
| 3  | 5   | om     | 运维群    | g      |
| 4  | 0   | po     | 运营群    | g      |
| 5  | 0   | od     | 研发群    | g      |
| 6  | 0   | sales  | 销售群    | g      |

#### casbin_object
| id | v             | 	name	 | created_at | 	updated_at	 | 
|----|---------------|--------|------------|--------------| 
| 1  | /api/v1/data1 | 数据源1   |            |              |
| 2  | source        | 源数据    |            |              |

#### casbin_policy

group 通过查询 casbin_user 和 casbin_role 进行自动生成
这里数据量小，act 直接使用字符串

| id | ptype | role     | 	objects | act     | created_at | 	updated_at	 |
|----|-------|----------|----------|---------|------------|--------------|
| 1  | p     | 2        | [1,2]    | read	   |  
| 2  | p     | 2        | [2]	     | write		 |
| 3  | p     | 3        | [2]      |         |

#### casbin_user
| uid | roles | 	  created_at | 	updated_at	 | 
|-----|-------|---------------|--------------| 
| 1   | [1,2] |               |              | 
| 2   | [5]   |               |              |  

