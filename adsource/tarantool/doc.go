// docker run --rm -it --label service.name=tarantool --label service.port=3301 tarantool/tarantool:2.2.0

package tarantool

/*
Create the advertisement space

```sh
cmp = box.schema.space.create('campaigns')
cmp:format({
	{name = 'id', type = 'unsigned'},
	{name = 'band_name', type = 'string'},
	{name = 'year', type = 'unsigned'}
})

box.schema.user.grant('guest', 'read,write,execute', 'universe')
```

*/
