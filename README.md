# ont_doc_update
用来更新ontology文档
	
## 名词解释

1. 源文件

    待添加到文档中心里的文件，分布在各个项目库中，包括文档库。

2. 文档文件

    文档中心的文档，与源文件会略有不同，主要是相对链接不同，因为文档中心里各个文档的路径和文件名都变了。

## Ontology 文档中心结构

`https://ontio.github.io/documentation/`被称为文档中心，`https://github.com/ontio/documentation`为文档库；文档中心的内
容来自于文档库下的`docs`目录，任何只在文档库或其他地方修改的内容，只要没有同步到`docs`目录下，都不会在文档中心里生效。

- `docs/lib/images/`
	
	该目录存放文档中用到的图片

- `docs/lib/_data/sidebars`
	
	该目录存放文档中心的组织结构
	
- `docs/lib/pages`
	
	该目录存放文档中心里各文档的内容

## 本项目配置文件说明

1. <b>`doc-map.json`</b>

    此文件为json格式的配置文件，文件内容为数组格式，存储的是所有源文件与文档文件之间的映射关系，每条映射中的"OriginalLink"为
    源文件下载链接，"NewPostion"为文档文件在文档库中的相对路径。程序运行过程为读取`doc-map.json`，下载"OriginalLink"指向的文件，
    做内容处理，将处理后的内容输出到"NewPostion"，所以`doc-map.json`指明了源文件与文档文件的对应关系。
    
    这是一个映射的例子：
    ```
    {"OriginalLink":"https://raw.githubusercontent.com/ontio/ontology-smartcontract/master/smart-contract-tutorial/SmartX_Tutorial.md",
        "NewPostion":"docs/pages/doc_en/Dapp/SmartX_Tutorial_en.md"}
    ```
    OriginalLink为源文件下载链接，注意该链接的主机名为`raw.githubusercontent.com`，不是`github.com`；
    
    NewPostion为文档文件在文档库中的相对路径。
    
    根据这条映射，本程序会将`ontology-smartcontract`库中`master`分支下的`smart-contract-tutorial`目录下的`SmartX_Tutorial.md`转换为文档库
    中`docs/pages/doc_en/Dapp/`目录下的`SmartX_Tutorial_en.md`。另一方面，该文件在文档中心的访问链接为
    `https://ontio.github.io/documentation/SmartX_Tutorial_en.html`，注意该链接的后缀为`.html`，不再是`.md`了。

2. <b>`link-map.json`</b>

    此文件为json格式的配置文件，文件内容为Map格式，存储的是源文件中的相对链接与对应的文档文件中的链接的映射关系，在处理过程中，源文件中的相对
    链接将会被替换成映射到的链接。
    例如：
    
    ```
    "documentation/get_started_cn.md[English](./get_started_en.md)":
    "./tutorial_for_developer_en.html",
    ```
    
    这个映射，第一行为源文件中的相对链接，源文件为`documentation`库中的`get_started_cn.md`，原链接为该文件中的`[English](./get_started_en.md)`
    ，这说明在源文件中，点击`English`将会跳转到`./get_started_en.md`，与之对应的文档文件中的链接则为`./tutorial_for_developer_en.html`，
    这代表在文档文件中点击`English`将会跳转到`./tutorial_for_developer_en.html`。
    
    处理时，该文件内容会自动检测源文件中的相对链接，如果该链接在`link-map.json`中不存在，则会将该链接添加到`link-map.json`中，并且该链接
    对应的映射链接为空。就像以下这种结构：
    
    ```
    "documentation/get_started_cn.md[English](./get_started_en.md)":
    "",
    ```
    
    这种设计能帮我们自动找出源文件中未被处理处理的相对链接，这些链接找出之后，我们需要手动添加链接映射，然后重新运行程序，这样源文件才会变成
    正确的文档文件。

## 本程序做的事情

   源文件的内容不一定能适应文档中心的结构，所以本程序做的就是将源文件转换成文档文件。过程如下：

1. 根据`doc-map.json`的"OriginalLink"下载源文件，做如下的内容处理：

	1. 检查源文件，取出相对链接，存放到`link-map.json`中，改成适应于文档中心的链接；

	2. 检查源文件标题格式，如果不符合，修正；

	3. 给文档文件添加文件头。

2. 内容处理完成后，将新的内容写入到`doc-map.json`中对应的"NewPostion"文件中。

目前已经配置了24小时文档更新服务器，它所做的事如下：

1. 下载最新的文档库到服务器上；

2. 下载本项目库的最新版本到服务器上；

3. 运行本项目中的`main.go`，将源文件转换成文档文件，并存放在本地文档库中的`docs/lib/pages`内的对应目录里，即更新本地文档库；

4. 将本地文档库的更新推送到GitHub上；

5. 文档中心更新完成。

	
## 文档中心更新流程

由于已经配置了每24小时自动更新文档中心，所以无需手动运行本程序来更新文档中心；文档内容有更新或者有新增文档的，只需完成以下步骤，在24小时之内就会
自动同步到文档中心了。

1. 编写或者更新源文件；

2. 将文件中引用的图片上传到文档库中的`docs/lib/images/`目录下；

3. 按照前文[所述](#本项目配置文件说明)，直接在GitHub上修改本项目库中的`doc-map.json`和`link-map.json`，修改后的文件保存做commit即可；如果只是更新
已有的文档，则无需修改`doc-map.json`，如果新的文档不涉及相对链接的改动，则`link-map.json`也不需要改；

4. 新增文档时，如果需在文档中心添加侧边栏导航链接，则需要修改文档库内`docs/lib/_data/sidebars`目录下的对应文件，增加相应的导航栏。
