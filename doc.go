package netns

/*
 * Wrapper around the setns() syscall to switch your application in the network
 * namespace of another process. It also get your Rx/Tx stats of the virtual devices.
 *
 * The process calling `netns.Setns()` should be root otherwise, a permission denied error
 * will be returned.
 */
