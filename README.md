# retreive data from SH1ES136 solar inverter with SOFARSOLAR datalogger 

As I tried to understand the already existing python codes to retreive the data from the inverter, I've decided to rewrite te python code into golang code. Why golang? because my final goal is to add the code into docker containers. And to make as small as possible docker images, it is the easiest way for me to use golang code.

## Initial release

(*) retreive inverter data
(*) retreive inverter string data
() get HW data

## For test purpose I've also added a stub to test, as in my case in the evening the solarinverter is down and I was not able to test
