# Data Integration Platform

The project is under construction. I should add a proper readme documentation asap

Dip, is a platform designed for data integrations from different sources.
For the time being, the platform supports 2 kinds of mechanisms to integrate
new sources. The push mechanism and the pull mechanism and only for the http
protocol but the goal is to expand this to other protocols and technologies
as well.

The platform's functionalities are developed in a micro-services model

![](./images/arch_high_level.png)

### Push Mechanism

![](./images/push_event_receptor.png)

### Pull Mechanism

![](./images/pull_event_receptor.png)