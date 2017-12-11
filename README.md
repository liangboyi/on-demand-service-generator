# What you need
### service-release
### service-adapter
### service-on-demand-broker-sdk

# 1.Create a Service Release
you can create a new one, here is a simple-redis-example
```
cd on-demand-service-generator
cd example-release
./update-release.sh
```
# 2.Generat sample code
use generator.sh to init a service-adapter-sample
```
cd on-demand-service-generator
./generator.sh redis
```
you can find a redis folder here
```bash
➜  on-demand-service-generator git:(master) ✗ cd redis
➜  redis git:(master) ✗ ls -l
total 0
drwxr-xr-x 3 bliang 102 Dec 11 10:22 local-test
drwxr-xr-x 4 bliang 136 Dec 11 10:22 redis-on-demand-broker
drwxr-xr-x 6 bliang 204 Dec 11 10:22 redis-service-adapter
drwxr-xr-x 9 bliang 306 Dec 11 10:22 redis-service-adapter-release
```

# 3.create your service-adapter
```
cd redis/redis-service-adapter
#write the code here
```
test the code
```
cd on-demand-service-generator
cd redis/local-test/service-adapter-test
./test.sh
```


# 4. Create your service-adapter release
```
cd redis
cp -r redis-service-adapter/* redis-service-adapter-release/src/redis-service-adapter
```
# 5.upload service-on-demand-broker-sdk-release
1 Download [on-demand-services-sdk 0.18.0](https://network.pivotal.io/products/on-demand-services-sdk#/releases/8516)
2 Upload release
```
bosh -e vbox upload-release on-demand-service-broker-0.18.0.tgz 
```
```bash
➜  redis git:(master) ✗ bosh -e vbox releases
Using environment '192.168.50.6' as client 'admin'

Name                           Version             Commit Hash
kafka-example-service          0+dev.3*            non-git
kafka-example-service-adapter  0.15.0-rc-1+dev.1*  b01dea5
learn-bosh                     0+dev.2*            8ecf7b4+
on-demand-service-broker       0.18.0-alpha-1*     1322b05
redis                          0+dev.1*            non-git
self-service-adapter           0+dev.4*            non-git
~                              0+dev.3             non-git
~                              0+dev.2*            non-git
~                              0+dev.1             non-git

(*) Currently deployed
(+) Uncommitted changes

9 releases

Succeeded
➜  redis git:(master) ✗
```


