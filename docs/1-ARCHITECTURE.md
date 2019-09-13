# Architecture
The general idea is to split each distinct part of the application 
as a whole into smaller services. For example a single service 
could handle all the data related to the rules for a specific
game and another one keeps track and modifies tournaments.

Communication between services is done through passing messages
to `Beanstalkd` ([more details](https://beanstalkd.github.io/)). 
This service acts as a centralized work queue.

An example being sending an email notification to a user. To do
this you would queue up a message with the email details and then
the `notifications` service will poll work and actually send off
the email when it is ready to process the message.

Example of flow where a service, `ServiceA`, wants to send an
email notification to a user. `ServiceA` will send a json blob
to `beanstalkd` which will queue it up for future processing. 
Then the `notification` service recieves and process messages
from the queue, which is where the actual email sending happens.
this pricipel should apply for any kind of processing that is 
not trivial allowing frontend API:s to respond instantly, and
defer work to later.

```
                                                                               Waiting to process
                                                                                   messages
+---------------+                   +-----------------+                +-------------------------+
|               |   SendEmailMsg{}  |                 |  RecvMessage   |                         |
|   ServiceA    +------------------->   Beanstalkd    +--------------->+   NotificationService   |
|               |                   |                 |                |                         |
+---------------+                   +-----------------+                +------------+------------+
                                                                                    |
                                                                          SendEmail | 
                                                                                    |
                                                                                    v
```

