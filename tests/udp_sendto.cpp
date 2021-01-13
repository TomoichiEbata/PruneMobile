// g++ udp_sendto.cpp -o udp_sendto で大丈夫だが、g++ -std=c++11としないと動かない(ことがある)

#include <stdio.h>
#include <string.h>
#include "simple_udp.h"

simple_udp udp0("192.168.0.8",12345); // ホストOSのUDPの受信側

int main(int argc, char **argv){
  udp0.udp_send("hello!");
  return 0;  
}

