set SRC_DIR=..\protobufFiles\
set DST_DIR=..\protos\

protoc -I=%SRC_DIR% Test.proto --go_out %DST_DIR% Test.proto
pause