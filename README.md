# 简介

实现将[B站录播姬](https://rec.danmuji.org/)录制的XML弹幕转换为ASS字幕文件。

## 主要功能

- [x] XML文件转换为ASS字幕文件
- [x] 批量处理XML文件
- [x] 自定义样式
    - [x] 弹幕透明度
    - [x] 弹幕描边颜色
    - [x] 弹幕阴影颜色
    - [x] 是否设置粗体
    - [x] 弹幕描边粗细
    - [x] 弹幕阴影大小
- [x] 指定弹幕显示时间
- [x] 弹幕时间平移
- [x] 限制弹幕显示范围
- [x] 调整弹幕的上下间距
- [x] 限制同屏弹幕密度
- [x] 按关键词过滤弹幕
- [x] 转换弹幕样式（如：底部弹幕转换成滚动弹幕）

## 开始使用

### 下载

- [Github](https://github.com/Hami-Lemon/converter/releases)
- [阿里云盘](https://www.aliyundrive.com/s/3U2mpc6oeLj) 提取码:`9aa5`

### 设置

设置文件需要和程序位于**相同的目录**,并且文件名为`setting.json`。
一个基本的设置文件内容为：

```json
{
  "fontsize": 26,
  "fontName": "黑体",
  "alpha": 0.3,
  "outlineColor": {
    "rgb": "0x49516A",
    "alpha": 0.1
  },
  "shadowColor": {
    "rgb": "0x49516A",
    "alpha": 0.1
  },
  "rollTime": 15,
  "fixTime": 5,
  "timeShift": 0,
  "bold": true,
  "outline": 0,
  "shadow": 1,
  "width": 1920,
  "height": 1080,
  "rollRange": 1.0,
  "fixedRange": 1.0,
  "spacing": 0,
  "density": 0,
  "overlay": false,
  "keyword": null,
  "convert": "s -> r"
}
```

上面展示的值都是默认值，当不存在设置文件，或者设置文件中对应设置项没有设置时，则会使用默认值。
同时，设置文件的内容**严格区分大小写**。即：`fontsize`是有效设置，而`Fontsize`则是无效设置。
此外，还需要注意值的类型，错误的类型可能会出现错误。

### 设置项说明

1. `fontsize`: 弹幕的字体大小。
2. `fontName`: 弹幕所用的字体，此字体应该安装在系统中（如果不生效，可以尝试修改名称为字体对应的英文名）。
3. `alpha`: 弹幕透明度，范围：[0.0, 1.0]，值越大越透明，即：0为全不透明，1为全透明。
4. `outlineColor`: 弹幕描边颜色，rgb表示该颜色的rgb值，alpha表示透明度，范围[0.0, 1.0]。
5. `shadowColor`: 弹幕阴影颜色。
6. `rollTime`: 滚动弹幕的显示时间，单位：秒。
7. `fixTime`: 顶部弹幕和底部弹幕的显示时间，单位：秒。
8. `timeShift`: 弹幕时间平移，单位：秒，负数表示弹幕时间前移（如果弹幕平移后时间小于0，则时间会被设置0），正数表示时间后移。
9. `bold`: 弹幕是否是粗体。
10. `outline`: 弹幕描边大小。
11. `shadow`: 弹幕阴影大小。
12. `width` `height`: 对应视频的分辨率。
13. `rollRange`: 滚动弹幕的显示范围，[0.00, 1.00]。
14. `fixedRange`：顶部弹幕和底部弹幕的显示范围，[0.00, 1.00]。
15. `spacing`: 弹幕的上下间距（这个字的大小受所用字体影响，同样的值，使用不同的字体会有不同的显示效果）。
16. `density`: 同屏弹幕密度限制，0表示无限密度。
17. `overlay`: 是否允许弹幕重叠。
18. `keyword`: 一个字符串数组，用于根据关键词过滤弹幕。
19. `convert`: 弹幕样式转换规则。

#### keyword 项说明

这一项的值应该是一个字符串的数组，其内容将被用于根据关键词过滤弹幕。
例如：`"keyword": ["不好不好", "赢"]`,将会过滤掉包含`”不好不好“`,`"赢"`的弹幕。

#### convert 项说明

这项用于指定弹幕类型转换，可以通过其实现：过滤底部弹幕，转换底部弹幕为滚动弹幕等功能。

符号定义：

- s: SC
- r：滚动弹幕
- t：顶部弹幕
- b：底部弹幕
- _: 过滤掉对应项弹幕

语法：[s,r,t,b] -> [s,r,t,b,_]。

`->`作为分隔符，分隔原弹幕类型和转换后的弹幕类型，分隔前面的内容表示原弹幕类型，后面的内容则是转换后的弹幕类型。

注意：`[]`表示其内容可以是方框中的任意个元素，例如：`[s,r,t,b]`表示其内容可以是
`s`,`r`,`t`,`b`,`sr`,`st`,`sb`,...,`srtb`中的任意一个。
额外说明：原弹幕类型和其对应的目标弹幕类型应该和分隔符处于对称位置。

![格式说明](http://i0.hdslb.com/bfs/album/ffe458f4dff17de3dab137ba2b4cd8fc9fdd16d3.png)

例如：`s -> r`表示将SC转换成滚动弹幕，`sb -> r_`表示将SC过滤掉，底部弹幕转换成滚动弹幕。

### 使用方式

```bash
Bullet Chat Converter (bcc) version: 0.1.0
将B站录播姬录制的XML文件转换成ass文件。

用法：bcc -x xml
  -x string
        待转换的xml文件，如果该路径是一个目录，则处理目录下的所有xml文件，默认为当前目录
  -h
  		显示帮助信息
```

`-x`选项指定xml文件所在的路径，可以是一个具体的文件，也可以是一个目录，如果是目录，这会处理当前目录下的所有XML文件，默认是程序运行的目录。

例如：

```bash
bcc -x demo.xml
```

## 效果展示

#### 无限密度，不重叠

![image-20220403210538703](http://i0.hdslb.com/bfs/album/ff441042a991de75d2a596d5db070f1c7fa419e2.png)

### 同屏密度限制100，不重叠

![image](http://i0.hdslb.com/bfs/album/18bba8db4a5c78954ce074d5c6d2b7243822542a.png)

### 滚动弹幕显示区域30%，无限密度，不重叠

![image-20220403211042525](http://i0.hdslb.com/bfs/album/fc90aff8fa3e79027628274e9c6aa751555fe8ab.png)

### 无限密度，允许重叠

![image-20220403211145997](http://i0.hdslb.com/bfs/album/a7ad47401cc3ec794495e2c91506871a3b03fc4e.png)

### 40%显示区域，同屏密度120，不重叠，底部弹幕转为滚动弹幕

![image-20220403211846869](http://i0.hdslb.com/bfs/album/267ad04405d32a91897249d90ab7fe75073f7dd1.png)

## 源码编译

本项目使用Golang开发，配置好Golang环境后可直接编译。

```bash
git clone git@github.com:Hami-Lemon/converter.git
cd ./converter/bcc
go build .
```

