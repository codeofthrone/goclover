# goClover

base on two open source project
https://github.com/warmans/golocc
and 
https://github.com/deadcheat/cloverleaf

#### Usage
```bash 
// using golang test to generate coverprofile
$ go test -v -coverprofile=c.out -covermode=atomic -cover  ./...
// convert c.out to clover format
$ goclover -f c.out -o clover.xml

```

#### Output example
```xml
<coverage clover="" generated="1543820003">
    <project name="" timestamp="1543820003">
        <package name="dpqqa.com/qas-2/Backend/controller/dms">
            <metrics loc="1811" cloc="8" ncloc="1803" packages="1" struct="8" method="0" exportedmethod="0" methodloc="0" function="1" exportedfunction="91" functionloc="1446" import="113" ifstatement="78" switchStatement="14" test="0" benchmark="0" assertion="80"></metrics>
            <file name="dpqqa.com/qas-2/Backend/controller/dms/create.go" path="">
                <line num="14" type="stmt" count="0"></line>
                <line num="15" type="stmt" count="0"></line>
                <line num="16" type="stmt" count="0"></line>
                <line num="17" type="stmt" count="0"></line>
                ...
            </file>
        </package>    
    </project>
</coverage>

```

### limitations:
```

Source code must place at $GOPATH/src
```


##### Reference
- clover format
http://openclover.org/doc/manual/4.2.0/ant--clover-report.html
- golang ast 
https://golang.org/pkg/go/ast/
