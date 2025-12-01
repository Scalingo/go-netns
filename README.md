> [!CAUTION]
> **This library was deprecated and archived on the 1st of December 2025.**

# [DEPRECATED] Network Namespace with Go

Wrapper around the `setns()` syscall to switch your application in the network
namespace of another process.

The process calling `netns.Setns()` should be root otherwise, a permission denied error
will be returned.

## Testing

The tests require Docker running on the local workstation. They have to be run as
root because of the constraints around `setns()`
