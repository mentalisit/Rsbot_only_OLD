.PHONY: all
all:clean build archive

clean:
	del .\Rsbot_only.zip
build:
	go build
archive:
	7z a Rsbot_only.zip Rsbot_only.exe
push:
	git push