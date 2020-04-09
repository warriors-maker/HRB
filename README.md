# HRB
Hash + Coded Version of Reliable Broadcast

#### Note:
The branch called **withoutlog** is my original implementation, you can use that branch to compare your new algorithms with all the algorithms I implement before.

You should fork this branch called **template** for your own reliable broadcast implementation.


# General Overview of RMB:
## Config file:
**trusted**: represents all your trusted ip

**faulty**: represents all your fauly ip

**algorithm**: represents your algorithm nunmber

**source_byzantine**: indicates whether your source can equivocate or not.

**data size**: data sizes in byteLength

**rounds**: # of messages broadcast 

**Note**: 

1.Even though we specify our trusted and faulty nodes, our algorithm needs to itself figure out which Ips are faulty.
Specification of the Ips in these two fields are just for setup node's behaviour purposes.

2.The first ip in the trusted ip will always to be the source node.

3.If you put source_byzantine = true, the first ip in the Trusted node will be the faulty source nodes.


## RMB supports both local mode and cluster mode.

**local Mode (meanning locolhost)**:
run go run main.go [index], were index is the index of your the concatination of the trusted + faulty nodes list.

__For example__ 

if you have trusted = [xxx:8000, xxx:8001] and faulty = [xxxx:8002]

If you run go run main.go 2, this command starts running the faulty node.

**Important Note**:

Notice that since local mode run all the nodes inside the same machine, you need to be careful when picking ports.

Supposed you have a node having ip address and port with 127.0.0.1:8000, it will also uses 8500 and 9000 port.

**Remote Mode (meaning each node resides in different machine)**
run go run main.go

# How to implement your own algorithm in RMB:

You need to do two things:

**First**, Please follow the format in the package called **HRBAlgorithm**, basically whenever you want to have a new algorithm, 
you need to create a new package.

**Second**, You need to touch the class called **ServerMain.go** inside the **Server Package**.

## Explanation of HRBAlgorithm Package:
#### algorithm.go: contains the variable you need to use for your algorithm.
You need to expose three functions to be used by the **Server Package**:
1. **AlgorithmSetUp**: initialization of your algorithm's variable.
2. **FilterRecData**: a filter function, this filter function will trigger appropriate function based on your message type.
For example, if a message type is a type of echo, it triggers the echoHandler().
3. **SimpleBroadcast** broadcast function: this function determines what kind of message you want to broadcast.

## What you need to do in ServerMain.go:
In **ProtocalStart()**, whenever you add a new algorithm, assign a integer to it.

Following the format inside this function, based on the algorithm field inside the config file, it will run the appropriate 
broadcast algorithm.

Inside:

1. Initialize your algorithm **HRBAlgorithm.AlgorithmSetUp**
2. **filter()**: this filter function will keep reading data from the readChannel created in 2.1, 
  and then call your algorithm **FilterRecData** to filter the data based on the mesage you get from the readChannel.
3. If current node is a source, it should invoke your **SimpleBroadcast**.
  
 

