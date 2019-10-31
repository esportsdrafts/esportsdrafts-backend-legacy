# Security
Much secure.

## Trust Model
The Kubernetes cluster where the code is deployed is the sole trusted space. In
other words the only place where secrets can be read and used. This means that
to perform and upgrade of the platform it has to be intiated by the cluster
itself, not an outside entity.
