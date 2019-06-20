# README

## model.go
gorm修改时间hook后, `db.Exec()`还是出现异常, `db.Delete()`正常

```log
// driver: github.com/lib/pq
// db.Exec()
pq: invalid input syntax for integer: "2019-06-20 03:19:29.558847646Z" 

[2019-06-20 11:19:29]  [13.23ms]  update videos set deleted_at = '2019-06-20 03:19:29' where scenic_id = 39 and vod_id = any('{86330106002001}')

// db.Delete()
[2019-06-20 11:25:29]  [25.43ms]  UPDATE "videos" SET "deleted_at"=1561001129  WHERE "videos"."deleted_at" IS NULL AND ((scenic_id = 39 and vod_id = any('{86330106002001}')))  
[0 rows affected or returned ] 
```

**个人感觉使用xorm比gorm更友好**

