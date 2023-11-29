# Overview

A wrapper around the [Chi](https://github.com/go-chi/chi) web framework, adding strongly typed routing and responses in JSON, and a global error handler.

All while not breaking the `http.Handler` interface.

The end goal is for this to be swapped out with the std library http router, once that becomes available.
