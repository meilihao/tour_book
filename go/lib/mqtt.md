# mqtt
# paho.mqtt.golang

## Publish
### Publish()
1. c.Publish()
   1. (c *client) Publish()

      `c.obound <- &PacketAndToken{p: pub, t: token}`
   1. `(c *client) startCommsWorkers()`

        1. `commsobound <- msg`
        2. `startComms(c.conn, c, inboundFromStore, commsoboundP, commsobound)`
            1. `startOutgoingComms(conn, c, oboundp, obound, outboundFromIncoming)`

                1. `case pub, ok := <-obound`

### Publish() conn write失败触发重连(by 关闭nanomq)
1. startOutgoingComms()

    `errChan <- err`
    
    1. `startComms()`

        `outError <- err`

        1. `(c *client) startCommsWorkers()`

            1. `case err, ok := <-commsErrors`
            2. `(c *client) internalConnLost()`

                1. `c.stopCommsWorkers()`
                2. `go c.reconnect(reConnDone)`

                    1. `(c *client) attemptConnection()`

                        1. `(c *client) startCommsWorkers()`

                            `c.conn = conn`
## [keepAlive](https://blog.csdn.net/HONGcheng930728/article/details/111137375)

# paho.golang