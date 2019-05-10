# postgis
参考:
- [GIS基本概念](https://blog.csdn.net/alinshen/article/details/78503333)

## install
```
$ sudo apt install postgresql-11-postgis-2.5
```

> deb位置: http://apt.postgresql.org/pub/repos/apt/pool/main/p/postgis/

## 使用
```
CREATE TABLE IF NOT EXISTS scenic_resort
(
    id            serial,
    name          VARCHAR(200) NOT NULL,                            -- 景区名称
    geog_point    GEOMETRY(Point, 4326),                            -- 坐标点
    geog_area     GEOMETRY(MULTIPOLYGON, 4326),                     -- 坐标区域
    voice         INTEGER[]    NOT NULL                            --语音id
);

INSERT INTO scenic_resort(name, geog_point, geog_area,voice)
VALUES ('杭州西湖风景区2', ST_GeomFromText('POINT(120.437504 30.045546)', 4326),
        'SRID=4326;MULTIPOLYGON(((120.151344 30.260626,120.125252 30.240829,120.146452 30.228668, 120.163361 30.245093,120.151344 30.260626)))',
         '{1,2,3,4}');
```

```sql
test=# select ST_AsEWKB(geog_point) from scenic_resort;
                      st_asewkb
------------------------------------------------------
 \x0101000020e6100000a1f7c610001c5e40b07614e7a80b3e40
 \x0101000020e6100000a1f7c610001c5e40b07614e7a80b3e40
(2 行记录)

test=# select ST_AsEWKT(geog_point) from scenic_resort;
               st_asewkt
---------------------------------------
 SRID=4326;POINT(120.437504 30.045546)
 SRID=4326;POINT(120.437504 30.045546)
(2 行记录)
```

> WKT(Well-known text)是开放地理空间联盟OGC（Open GIS Consortium ）制定的一种文本标记语言，用于表示矢量几何对象、空间参照系统及空间参照系统之间的转换
> WKB(well-known binary) 是WKT的二进制表示形式，解决了WKT表达方式冗余的问题，便于传输和在数据库中存储相同的信息

## FAQ
### SRID 4326 26910区别
4326 : WGS的最新版本为WGS 84（也称作WGS 1984、EPSG:4326), GPS使用的坐标系, **推荐**
26910 : [NAD83 UTM Zone 10N (EPSG:26910)](https://epsg.io/26910), 仅支持在WGS84的[-126.0000, 34.4000, -120.0000, 77.0000]区间, 不推荐

### 无法打开扩展控制文件 "/usr/share/postgresql/11/extension/postgis.control": 没有那个文件或目录
未安装postgis

安装后在psql对应的数据库中执行`CREATE EXTENSION postgis;`即可启用postgis, 也可用`\dx;`检查已安装的扩展

### MULTIPOLYGON 和 Polygon
MULTIPOLYGON是Polygon

POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10)) # 实心
POLYGON ((35 10, 45 45, 15 40, 10 20, 35 10), (20 30, 35 35, 30 20, 20 30)) # 空心