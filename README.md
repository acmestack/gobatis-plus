# gobatis-plus


## Usage
```
gobatis-plus -i {YOUR_PACKAGE_DIR} -o {OUTPUT_DIR}
```
example:
```
gobatis-plus -i github.com/acmestack/gobatis-plus/example -o tmp
```

## Package
1. Add doc.go to enable gobatis in package
```
// +gobatis:enable
package example
```
2. Add your code
```
// +gobatis:data:tablename=tbl_user
type UserDo struct {
	// +gobatis:tableid:value=user_id,idType=auto
	UserId int64

	// +gobatis:tablefield:value=user_name
	UserName string

	// +gobatis:tablefield:value=status
	// +gobatis:tablelogic:value=0,delval=1
	Status int8

	// +gobatis:tablefield:value=create_time,fill=insert
	CreateTime time.Time

	// +gobatis:tablefield:value=rec_var
	// +gobatis:version
	RecVersion uint64
}

// +gobatis:mapper
type UserMapper interface {
	mapper.BaseMapper[UserDo]
	
	// +gobatis:select="select * from tbl_user where id = #{UserDo.UserId}"
	Select(user UserDo) (users []UserDo, err error)
}
```