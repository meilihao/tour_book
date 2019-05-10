# postgis
参考:
- [GIS基本概念](https://blog.csdn.net/alinshen/article/details/78503333)
- [GeoJSON对象](https://www.jianshu.com/p/5c6c6e76d4df)

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

test=# select ST_AsBinary(geog_point) from scenic_resort;
                 st_asbinary
----------------------------------------------
 \x0101000000a1f7c610001c5e40b07614e7a80b3e40


test=# select ST_AsEWKT(geog_point) from scenic_resort;
               st_asewkt
---------------------------------------
 SRID=4326;POINT(120.437504 30.045546)

test=# select ST_AsText(geog_point) from scenic_resort;
          st_astext
-----------------------------
 POINT(120.437504 30.045546)

 test=# SELECT ST_GeomFromText('POINT(1.2345 2.3456)');
              st_geomfromtext
--------------------------------------------
 0101000000A1F7C610001C5E40B07614E7A80B3E40
```


point:
```
EWKT:
SRID=4326;POINT(120.437504 30.045546)

EWKB:
01 // 编码序: 00为使用big-endian编码(XDR)，01为使用little-endian编码(NDR)
0100 // 矢量数据的类型, 0100代表Point
0020 // 矢量数据的维数, 0020代表该点是二位的
e6100000 // 矢量数据的空间参考SRID, E6100000是4326的整数十六位进制表达
a1f7c610001c5e40b07614e7a80b3e40 // structPoint
```

MultiPolygon:
```
EWKB:
SRID=4326;MULTIPOLYGON(((120.437504 40,20 45,45 30,120.437504 40)),((120.437504 35,10 30,10 10,30 5,45 20,120.437504 35),(120.437504 20,20 15,20 25,120.437504 20)))

EWKT:
01 // little-endian
0600 // MultiPolygon
0020 // (lng,lat)
e6100000 // SRID
02000000 // Polygons个数

01 03000000 01000000 04000000 // WKBPolygon结构 = 01(little-endian) + 03000000(wkbType) + 01000000(numRings) +04000000(numPoints)
a1f7c610001c5e400000000000004440 00000000000034400000000000804640 00000000008046400000000000003e40 a1f7c610001c5e400000000000004440 // Polygons1

01 03000000 02000000 06000000
a1f7c610001c5e400000000000804140 00000000000024400000000000003e40 00000000000024400000000000002440 0000000000003e400000000000001440 00000000008046400000000000003440 a1f7c610001c5e400000000000804140

04000000
a1f7c610001c5e400000000000003440 00000000000034400000000000002e40 00000000000034400000000000003940 a1f7c610001c5e400000000000003440
```

> WKT(Well-known text)是开放地理空间联盟OGC（Open GIS Consortium ）制定的一种文本标记语言，用于表示矢量几何对象、空间参照系统及空间参照系统之间的转换
> WKB(well-known binary) 是WKT的二进制表示形式，解决了WKT表达方式冗余的问题，便于传输和在数据库中存储相同的信息
> [WKB描述的几何对象](https://blog.csdn.net/yaoxiaochuang/article/details/53117693)
> [矢量空间数据格式(struct define)](https://www.cnblogs.com/marsprj/archive/2013/02/08/2909452.html)

## FAQ
### SRID 4326 26910区别
4326 : WGS的最新版本为WGS 84（也称作WGS 1984、EPSG:4326), GPS使用的坐标系, **推荐**
26910 : [NAD83 UTM Zone 10N (EPSG:26910)](https://epsg.io/26910), 仅支持在WGS84的[-126.0000, 34.4000, -120.0000, 77.0000]区间, 不推荐

### 无法打开扩展控制文件 "/usr/share/postgresql/11/extension/postgis.control": 没有那个文件或目录
未安装postgis

安装后在psql对应的数据库中执行`CREATE EXTENSION postgis;`即可启用postgis, 也可用`\dx;`检查已安装的扩展

### MultiPolygon 和 Polygon
MultiPolygon包含多个Polygon.
Polygon包含多个LinearRing(线性环).

POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10)) # 实心
POLYGON ((35 10, 45 45, 15 40, 10 20, 35 10), (20 30, 35 35, 30 20, 20 30)) # 空心
MultiPolygon (((40 40, 20 45, 45 30, 40 40)),((20 35, 10 30, 10 10, 30 5, 45 20, 20 35),(30 20, 20 15, 20 25, 30 20))) # 一个实心多边形+一个空心多边形