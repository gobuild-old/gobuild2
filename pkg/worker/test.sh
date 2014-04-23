redis-cli -r 3 RPUSH resque:queue:hello '{"class":"Hello","args":["hi","there"]}'
./worker -interval=1 -queues=hello -concurrency=1
