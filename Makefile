# 定义变量
BINARY_NAME = tool
DOCKER_IMAGE_NAME = yantao/$(BINARY_NAME)
DOCKER_TAG = latest

# 默认目标
all: build-image push-image

# 构建Go程序
build:
	go build -o $(BINARY_NAME) main.go

# 构建Docker镜像
build-image: build
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

# 推送Docker镜像到仓库
push-image: build-image
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# 清理生成的文件
clean:
	rm -f $(BINARY_NAME)