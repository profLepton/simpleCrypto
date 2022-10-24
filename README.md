# simpleCrypto
 ## Making cryptotechnology for everyone

 This project aims to increase the utility of crypto currency by devising new techniques that can take advantage of the blockchain.  

 # Messaging service protocol

This protocol, when adopted by enough nodes establishes a completely decentralized messaging service that is entirely operated through a peer to peer network.

The first version of the protocol uses the blockchain for payments only, and establishes a simple service system. The second version will use Multi signature contracts to further strengthen the service.

# How it works (current implementation)

Let us consider three parties. The sender (Alice), the server (Bob) and the receiver (Carol). When Alice firsts decides to communicate with Carol, Alice exchanges a series of encrypt only keys with Carol. Alice will encrypt messages with the keys from Carol, so that an intermediary may not read them (End to End encryption). 

When Alice is ready to send a message to Carol, Alice contacts any available server in her network, and we shall refer to the server that responds to Alice's request as Bob. 
Alice gives the encrypted message along with its intended recipient and "delivery fee" to Bob. "delivery fee" refers to the fee that Bob will collect upon delivery of Alice's message to Carol. Bob, if he thinks the fee is reasonable, will accept the message and encrypt it himself using a secret key and store it.

When Carol is ready to receive the message, Carol contacts the servers in her network if there are any messages for her. This request for messages to her will propagate in the network as each node that is asked for a messsage check will ask other nodes. If Bob is in her network, Bob will respond to Carol's request. 

Carol can verify that the message was indeed from Alice by verifying her siganture. Carol cannot read the message as it has been encrypted by Bob. To decrypt the message Carol requires the secret key which Bob used to encrypt the message. Bob releases this key upon receiving the payment from Carol. 

That is the underlying methodology.

## Current implementation

Our current implementation includes a 'server', which accepts new blocks from a 'client' and adds them to the block chain. A 'client' which receives transactions from programs names 'Alice', 'Bob' and 'Carol' and mines blocks with the transactions inside them.
'Alice' and 'Carol' send and receive messages, while Bob is the server that acts as a node to facilitate the communication.

## Future implementations

Future implementations can use Multi signature accounts or similar techniques to achieve a more reliable method of payment, allowing payment from the sender, or allow the sender and receiver to open an account with the server to receive and make continued payments, similar to the implementation of a lightning network.
