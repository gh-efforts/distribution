# distribution
> 用于filplus数据分配使用，默认使用添加数据集时设置的副本数量，不允许数据重复  

## 使用方法
### 增删改查用户/组织
#### 增加
```bash
$ ./dist user add -h
NAME:
   dist user add - add user

USAGE:
   dist user add [command options] [arguments...]

OPTIONS:
   --sp value, -s value   specify sp list
   --org value, -u value  specify org name
   --force                force update user,cover (default: false)
   --help, -h             show help
```
#### 删除
```bash
$ ./dist user  delete -h
NAME:
   dist user delete - delete a org

USAGE:
   dist user delete [command options] [arguments...]

OPTIONS:
   --org value, -u value  specify org
   --really-do-it         must be specified for the action to take effect (default: false)
   --help, -h             show help
```
#### 查看
```bash
$ ./dist user view
+--------+-----------------------------------------------------------------------------+
|  org   |                                     sps                                     |
+--------+-----------------------------------------------------------------------------+
|  beck  |                             ["f01001","f01002"]                             |
| beck-2 |                             ["f01001","f01002"]                             |
| beck-3 |                             ["f01001","f01002"]                             |
| beck-4 |                             ["f01001","f01002"]                             |
| beck-5 |                             ["f01001","f01002"]                             |
| beck-6 | ["f01001","f01002","13","53","asdf","fasd","","gasdf","","","gasdf","asdf"] |
+--------+-----------------------------------------------------------------------------+
```
### 增删改查数据集
#### 增加/更新数据集
```bash
$ ./dist dataset add -h
NAME:
   dist dataset add - add a dataset

USAGE:
   dist dataset add [command options] [arguments...]

OPTIONS:
   --name value, -n value       specify dataSet name
   --duplicate value, -d value  specify dataSet duplicate (default: 0)
   --filepath value, -f value   specify dataSet filepath. must include pieceCid,pieceSize,carSize
   --force                      force update dataset,cover (default: false)
   --help, -h                   show help

$ ./dist dataset add -n hofe -f test/dataset.json -d 10
```

#### 删除数据集
```bash
$ ./dist dataset delete -h
NAME:
   dist dataset delete - delete a dataset

USAGE:
   dist dataset delete [command options] [arguments...]

OPTIONS:
   --name value, -n value  specify dataSet name
   --really-do-it          must be specified for the action to take effect (default: false)
   --help, -h              show help
```
#### 查看数据集
```bash
$ ./dist dataset view 
+-------------+-----------+-------+----------+----------------+---------------------+
| dataSetName | duplicate | spSum | pieceSum | pieceSize(TiB) |    carSize(TiB)     |
+-------------+-----------+-------+----------+----------------+---------------------+
|    wyth     |     0     |   1   |    10    |     0.3125     | 0.17188961445026507 |
|    hofe     |     0     |   0   |    10    |     0.3125     | 0.17188961445026507 |
|    hofe1    |     0     |   0   |    10    |     0.3125     | 0.17188961445026507 |
|    hofe2    |     0     |   0   |    10    |     0.3125     | 0.17188961445026507 |
|   ho3fe2    |     0     |   0   |    10    |     0.3125     | 0.17188961445026507 |
+-------------+-----------+-------+----------+----------------+---------------------+
```
### 获取合适的文件下载链接
```bash
$ ./dist dataset get -h
NAME:
   dist dataset get - get the download link for the dataset

USAGE:
   dist dataset get [command options] [arguments...]

OPTIONS:
   --name value       specify dataSet name
   --sp value         specify a sp
   --size value       specify total pieceSize(TiB) (default: 0)
   --duplicate value  specify dataset duplicate (default: 10)
   --repeat value     specify dataset repeat (default: 0)
   --prefix value     specify url prefix [$DIST_PREFIX]
   --suffix value     specify url suffix (default: ".car") [$DIST_SUFFIX]
   --really-do-it     must be specified for the action to take effect (default: false)
   --help, -h         show help
   
# 先不要使用 --really-do-it , 确认无误后使用
$ ./dist dataset get --name hofe --sp f01001 --size 0.0625
$ ./dist dataset get --name hofe --sp f01002 --size 0.0625 --really-do-it
baga6ea4seaqinzfzpn2yzgshx4zduxwnk2sxwyu7uahyagn6swqwji66pvxriii.car
baga6ea4seaqagfkxwkfdmwskt7mw3hbglwad2ermgav766yxeybux3czntpfify.car
total pieceSize:0.0625, total carSize: 0.034377923078864114, missing pieceSize:0

```
### 增删改查piece
#### 增加/修改
> 修改需填写全部sps
```bash
$ ./dist piece add --name wyth --pieceCid baga6ea4seaqinzfzpn2yzgshx4zduxwnk2sxwyu7uahyagn6swqwji66pvxriii --pieceSize 34359738368 --carSize 18899463547 --sps f01002 --force
add piece baga6ea4seaqinzfzpn2yzgshx4zduxwnk2sxwyu7uahyagn6swqwji66pvxriii success!
```
#### 删除
```bash
$ ./dist piece delete -h                                                                                                                                                     
NAME:
   dist piece delete - delete piece

USAGE:
   dist piece delete [command options] [arguments...]

OPTIONS:
   --name value      specify dataSet name
   --pieceCid value  specify pieceCid
   --really-do-it    must be specified for the action to take effect (default: false)
   --help, -h        show help
```
#### 查看
```bash
$ ./dist piece view --name hofe
+------------------------------------------------------------------+----------------+--------------------+------+
|                             pieceCid                             | pieceSize(GiB) |    carSize(GiB)    | sps  |
+------------------------------------------------------------------+----------------+--------------------+------+
| baga6ea4seaqinzfzpn2yzgshx4zduxwnk2sxwyu7uahyagn6swqwji66pvxriii |       32       | 17.601497049443424 | null |
| baga6ea4seaqagfkxwkfdmwskt7mw3hbglwad2ermgav766yxeybux3czntpfify |       32       | 17.60149618331343  | null |
| baga6ea4seaqo5x5opegkbw2mzlz5jw4ckvphv7pkjncjiyuobxmfshoyduxaudy |       32       | 17.601496373303235 | null |
| baga6ea4seaqhsexpcqi3wkrrdkgatfga3wvwgibra37sdcmiqjpajayzm3fksiq |       32       | 17.60149618331343  | null |
| baga6ea4seaqaikuohit2k3jovndkl2bycqulju5g6ihfgqjpkhh237pxnpsc4ga |       32       | 17.601497495546937 | null |
| baga6ea4seaqoueiwop6ee4j7ziz5jj7m54m2quotyvrr56ylubckkivhrtz6giq |       32       | 17.601496262475848 | null |
| baga6ea4seaqazg5awsaupo3hyp7ro6xjeoagfp46e75tqakrxu3xnbfvtjcykla |       32       | 17.601496363058686 | null |
| baga6ea4seaqps2qhokkvqf42mrp7scqwmzn3tyibkwnxoe34twwjpssb74qk4ji |       32       | 17.601496262475848 | null |
| baga6ea4seaqefv5fpl546dmh3qvdkg7ukiamiybwhezg5m7ieqvts2fecqazsmy |       32       | 17.60149637144059  | null |
| baga6ea4seaqgs4yqakww6p4kihec7246s7knmgeaajqj5ehfnqzo6mhqgo3uybq |       32       | 17.601496652700007 | null |
+------------------------------------------------------------------+----------------+--------------------+------+
```