### 前置文档

```html
代码大全2
https://12factor.net/
https://github.com/golang/go/wiki/CodeReviewComments
http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
https://en.wikipedia.org/wiki/SOLID
```

### 开发者

#### 工具
```text
    尽量使用工具来完成书写格式上的规范性约束
    譬如使用 gofmt 来完成代码的格式化（或使用 goimports来替代 gofmt）
    强制执行 gofmt, go vet, go lint

    #例如，使用 vim 时，可以安装使用 vim-go插件，开启如下自动处理：
    let g:syntastic_go_checkers = ['gometalinter', 'govet', 'errcheck']
    let g:syntastic_mode_map = { 'mode': 'active', 'passive_filetypes': ['go'] }
    let g:go_metalinter_autosave = 1
    let g:go_fmt_command = "goimports"
    使用 vscode/goland 时，也有对应的插件和配置
    提交CR的代码需要使用 gofmt 或 goimports 进行格式化，统一格式
```

#### 命名及注释
```text
    在编码阶段同步写好变量、函数、包注释，注释可以通过 godoc 导出生成文档
    注释必须是完整的句子，以需要注释的内容作为开头，句点作为结尾
    程序中每一个被导出的（大写的）名字，都应该有一个文档注释
    （https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
    https://golang.org/doc/effective_go.html#commentary）
    命名最好是自注释的（请学习 《代码大全》）
    采用驼峰方式命名
    需要注释来补充的命名就不算是好命名
    使用可搜索的名称：单字母名称和数字常量很难从一大堆文字中搜索出来。单字母名称仅适用于短方法中的本地变量，名称长短应与其作用域相对应。若变量或常量可能在代码中多处使用，则应赋其以便于搜索的名称。
    做有意义的区分：Product和ProductInfo和ProductData没有区别，NameString和Name没有区别，要区分名称，就要以读者能鉴别不同之处的方式来区分 。
    函数命名规则：驼峰式命名，名字可以长但是得把功能，必要的参数描述清楚，函数名名应当是动词或动词短语，如postPayment、deletePage、save。并依Javabean标准加上get、set、is前缀。例如：xxx + With + 需要的参数名 + And + 需要的参数名 + …..
    结构体命名规则：结构体名应该是名词或名词短语，如Custome、WikiPage、Account、AddressParser，避免使用Manager、Processor、Data、Info、这样的类名，类名不应当是动词。
    包名命名规则：包名应该为小写单词，不要使用下划线或者混合大小写。
    接口命名规则：单个函数的接口名以”er”作为后缀，如Reader,Writer。接口的实现则去掉“er”。
    避免出现大量注释掉的代码或无用代码、注释，不需要的请直接删掉
    一行代码不要过长，适当使用换行来增加代码可读性
    一个函数避免过长，过长时考虑拆分或提取子函数
    一个文件避免过长，按照函数进行文件分类
```
    

#### 函数及语句
DRY
```text
尽量将功能代码抽取成独立的函数，并对函数进行单测
强制map 使用前要初始化
map 在并发场景中需要使用锁进行保护建议用 sync.Mutex 的指针类型来使用变量值，从而避免 Mutex 的值拷贝
在定义函数或方法时，如果该函数或方法有远程访问调用，需要增加一个 context 参数，且作为第一个传入参数
可以用于超时处理、传递内含数据（用户信息、trace信息等）
避免使用 interface{} (empty interface)作为函数传入参数，可是使用 Named Interface
```


声明  slice 变量时，只声明不定义
```go
var someSlice []string
// not
someSlice := []string{}
```
if 判断时，如果判断的变量为 boolean 类型时，不需要再与 true/false 进行对比:
```go
//不建议
if isPrefix == true {
    // blabla
    }
//建议
if isPrefix {
    // blabla
}
```
要明白函数返回值 return nil 与 return someVar (var someVar *someStruct == nil) 的区别
```go
type InternalError struct {
        msg string
    }

    func (ie InternalError) Error() string {
        return ie.msg
    }

    func main() {
        fmt.Printf("%#v\n", func() error { return nil }())
        fmt.Printf("%#v\n", func() error { var err *InternalError; return err }())
    }
```
```text
避免代码中出现 magic number or magic string,给每个 magic 值定义一个 常量
返回值中如果有struct，最好是返回指针类型的数据 如 func() (*someStruct, error)
返回值中如果有struct的slice，则可以直接返回改 struct 的 slice 如 func() ([]someStruct, error)
如果有变量引用了slice 或数组中的非指针数据，最好做一个临时拷贝再引用，会涉及到垃圾回收的效率。如：
```

```go
    for idx := range listPods.Items {
       if listPods.Items[idx].Metadata.Name == destPodName {
          // not: return &listPods.Items[idx]
          tmp := listPods.Items[idx]
          return &tmp
       }
     }
```
time.Ticker must be stopped
defer timer.Stop()

#### 错误处理
```text
函数返回值中有 error的，一定要处理（最差打一条日志）
优先处理错误，遇到错误可以提前处理，避免 if 嵌套
```

````go
func doSomething() error {
      if err := do1(); err == nil {
          if err = do2(); err == nil {
             return nil
          } else {
            // print err
            return err
          }
      } else {
         // print err
         return err
      }
}
func doSomething2() error {
      if err := do1(); err != nil {
         // print err
         return err
      }
      if err = do2(); err != nil {
         // print err
         return err
      }
      return nil
}
````
#### 设计与实现
```text
接口与实现分离
函数或类型要单一职责,学习一下：SOLID设计原则，https://en.wikipedia.org/wiki/SOLID）
注重实现效率，避免在一个for 循环内部进行多次 网络请求，可以一次请求完成后，再进行for循环处理
关注自己实现接口的耗时（设定自己的SLO）
"TODO" items typically have "TODO(yourname)" or "TODO(someone)" or "TODO(sig)"
```

#### 测试及CR
```text
开发新接口的，要提供接口测试case
尽量多写单测
```

```go
test table
mock interface(dependency injection,SOLID)
```
```text
在个人分支上开发时，适当对自己提交的同一个功能的多次 commit 进行压缩，保证提交 CR 时，只留一个 commit,如果 CR 已过，
后续有bug 需要修复时，可以再次提交修复的commit
CR 的代码越短越好（保证实现功能）
在一次 CR 提交过程中，尽量保证只有一个 功能进行 CR, 即使顺手改了一个其他的小功能，也分开提CR,
（使用多个 个人开发分支，每个分支上进行一个功能的修改，分开提 CR）
譬如 使用 feature_dev_reboot_1 进行一个新功能的开发，使用 feature_dev_reboot_2 进行另一个功能的bug 修复之类的使用方式
选取合适的reviewer
找一个功能相关者（backup）进行主功能 review
其他人也可以进行review
重要逻辑不少于 2 位的 reviewer
提前沟通好功能（文档、当面沟通）
提交 CR 之前自己先 Review 一遍 ,保证程序编译通过，功能测试通过
```


#### 评审者
- checklist , 清晰度
```text
能否看懂代码到底实现了什么(关键部分是否有相应的注释）
看不懂的地方就要 comment your questions
是否有相关的关键性注释文档（关键逻辑、功能、算法）
变量命名是否规范,函数命名是否规范,文件命名是否规范
是否有能提取的函数块
代码是否易于维护（简单性，复杂的实习肯定不是一个好的实现）
未来是否易于扩展或重构
关键步骤是否有日志打印
接口、远程调用是否有必要新增 metric 数据
代码中能否提前判断错误进而提前返回
错误字符串 不能以 大写开头、不能以标点符号结尾
是否有不必要的垃圾代码、注释
是否有被注释掉的代码（应该被删掉）
变量是否被定义在最小的作用于范围内
```


- 功能逻辑
```text
代码逻辑是否正确
算法是否正确
是否实现目标功能、目标功能是否有遗漏
修复是否必须（是否有已经存在的相关功能）
关键逻辑是否有单元测试
是否有相关的接口测试case
接口方法是否正确（GET、POST、PUT..）
接口是否有权限认证
接口发生变化时是否能够向前兼容
返回错误是否被忽略（resut, _ := someFunc()）
依赖是否是通过接口的方式被注入依赖方
goroutine 中是否有未被保护的共享数据
goroutine 是否有可能泄漏
channel 是否有可能被阻塞
sync.Mutex 是否有可能被拷贝
使用自定义的 http.Request 进行HTTP 请求时，在错误处理之后，需要 defer req.Body.Close()
使用自定义的 http.Request 进行HTTP 请求时，需要增加合适的 timeout
```


#### 性能、效率与安全
```text
输入参数是否被验证（validate the input）
相关逻辑是否会有潜在的性能缺陷（多余的网络交互）
相关算法是否有优化的空间
是否有越界行为
是否有未受保护的goroutine 共享数据访问行文
map 是否被初始化
代码是否易于扩展或修改
代码是有能被重用
错误是否被正确处理
是有有可能使用未正确初始化的指针变量
```


#### 其他
```text
要了解 CR 既是 Review的过程，也是学习过程
只对代码进行 comment 及 question
提问时要阐明自己的疑问，让 developer 看懂问题
```


