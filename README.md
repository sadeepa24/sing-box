# sing-box Mod

## Overview

This project extends the functionality of the `sing-box` universal proxy platform, adding real-time user management and enhanced features for dynamic configurations.

## Features

- **Same User Multiple Inbound Support**: Use the same configuration for multiple inbound.
- **Real-Time Inbound Changes**: Dynamically update user inbound settings without restarting.
- **Real-Time Outbound Changes**: Dynamically update user outbound settings without restarting.
- **Real-Time Usage Monitoring**: Track user connection status in real time.
- **Connection Limiting**: Apply IP and bandwidth limits.
- **Connection Management**: Close all connections for a specific user.
- **Comprehensive User Status**: Retrieve the status of all users at once.

Already Sing Box supports adding/deleting new inbound/outbound.

## New API Methods (of box.Box)

### User Management

- `AddUser(u opts.User) (opts.UserStatus, error)`: Add a new user.
- `AddUserReset(u opts.User) (opts.UserStatus, error)`: reset alredy added user status with new u.
- `RemoveUser(u opts.User) (opts.UserStatus, error)`: Remove an existing user.
- `GetStatusUser(u opts.User) (opts.UserStatus, error)`: Get the status of a specific user.

### Inbound/Outbound Management

- `ResetInbound(u opts.User)`: Reset inbound settings for a user according to new u.
- `ChangeOutbound(u opts.User) error`: Change outbound settings for a user according to new u.

### Connection Management

- `CloseAllConn(u opts.User)`: Close all active connections for a user.
- `AllUserStatus() map[int]opts.UserStatus`: Retrieve the status of all users.

## Official Documentation

https://sing-box.sagernet.org

The universal proxy platform.
