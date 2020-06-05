# extjs
## FAQ
### 如何查找table用到的data
1. 找到table渲染时用到的那个`columns`

    它会使用`dataIndex`属性作为用到的数据对象的key.
    部分字段可能会使用自定义渲染函数: `renderer`, 该函数的第一个参数即为`data["dataIndex"]`
1. 找`store = Ext.create('base'`, 通常在同名文件中.

    columns用到的数据对象来源于store
1. 找到创建store使用的`base`的model, response data转成store用到的数据对象的映射就在该model中.

    ```js
    Ext.define('my.model.StorageNode', {
        extend: 'Ext.data.Model',
        fields: [
        {
            name:  'displayname',
            convert: function (v, record) { // record为response data中相应数组中的成员
                var res = '';
                if (record.raw.doc.metadata != undefined)
                    res = record.raw.doc.metadata.displayName;
                return res;
            }
        },...]
    });
    ```