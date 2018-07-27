# ont_doc_update
用来更新ontology文档

1. 使用命令`go get -u github.com/qiluge/ont_doc_update`下载本仓库代码;

2. 使用命令`go get -u github.com/ontio/documentation`下载ontology document库;

3. `doc-map.json`存放的是原始文档和该文档在document库中未知的映射关系，如果新增文档，请在该文件中增加一行映射;

4. `link-map.json`存放的是原始文档中的相对链接与文档中心里对应文档中的链接的映射关系；

5. 运行main.go，会自动将原始文档内容做处理，然后放到document本地库的对应路径，并且将所有在`link-map.json`文件中配置的相对链接映射成文档中心的链接，并且新文档中的相对链接会被添加到该文件中;

6. 如果`link-map.json`文件中没有新增的行，则document更新完成，可以将document直接推到GitHub上，文档中心的内容相应更新；

7. 如果`link-map.json`文件中有新增的行，则需要补齐该相对链接对应的文档中心里的链接，然后再次运行main.go更新本地document库，
再将修改后的内容推到GitHub上。
