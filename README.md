# Tool đăng kí môn HUFLIT
## Cách sử dụng
1. Cài đặt docker
2. Build docker image: `docker build -t tool-huflit .`
3. Tạo file **input.txt** theo cú pháp: `<code>|<code LT>|<code TH>(nếu có, không có thì để trống)|<tên môn>`
4. Tiến hành chạy: `docker run -it --rm -v "$(PWD)/input.txt:/app/input.txt" tool-huflit ./tool -username <username> -password <password> -workers 10 -file input.txt`

