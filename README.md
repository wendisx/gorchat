> ## GORCHAT

```
My original intention was to make a small communication system,
but I don't have time to think about how to make this project 
go further, nor is there any need to make it go further. At first, 
I conceived a simple and efficient communication application, 
which could be a very simple program. Of course, this cannot be 
compared with real commercial products, and I know this very well. 
I may do it in the future, but reality is often not satisfactory. 
I need to combine a lot of technical points that look advanced, 
but in fact they are just some technical points that can be thought 
of by careful thinking to enrich my "project", but in fact this is 
not of much value. We don't have to deify an advanced technology. 
Excessive deification will only lead to ourselves forever regarding 
it as a kind of magic. Just remember one thing: technology comes from 
reality. The ideas behind it are what we should understand and learn 
in depth.
```

```
Project TimeLine

 .
 |
[*] --> basic backend framework(finished 07/10)
 |     |__ user section
 |     |__ single section
 |     |__ group section
 |
[ ] --> basic im system
 |
[ ] --> optimize
 |
 |
[ ] --> basic frontend
 |
 V
[X]

Project Tech

 -- Mysql 
 -- Redis/Valkey
 -- Kafka
 -- Docker
```

> ## Reference

- [bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [im system design](https://help.aliyun.com/zh/tablestore/use-cases/message-system-in-the-modern-im-system/?spm=a2c4g.11186623.help-menu-27278.d_5_4_3.14f5251ca1efmC)

```
Just study the above architecture and you will gain something from 
it. In fact, the entire backend operation logic is just so much. 
of course, this is just the basic part.
```

```
-- Date: 7/22
-- Author: wens
It took me some time to complete these basic parts, and I feel ashamed of it.
I think it is not necessary to spend too much time on these basic things, 
but this is indeed the first time I have completely implemented a chat backend 
from scratch. I think the most worthwhile thing during this period is that
I was able to find the most efficient development design, and I am very glad 
that I got it. Of course, this is a relative result based on continuous wrong 
designs. Maybe this can help others, but I don't care about it.Golang has a 
different model from Java. But the core is the same. Before developing anything,
I can think of the need for an early warning system, at least I can know which 
naughty process has a problem at the first time. I fell too much into the trap 
of structured early warning. Until all this happened, I looked at a lot of noisy 
text. I tried to understand everything that the early warning system tried to tell 
me, but it was just a conspiracy to try to separate me from my program. I failed.
The early warning system should tell me what happened as simply as possible while 
everything is working normally, instead of telling me some details. The structured 
early warning system can combine logs with control panels for visualization.The 
early warning system needs to be in a certain format, at least the warning level 
is an indispensable part. This is more like a signal-based early warning system. 
Some warnings are just unnecessary or can be used as warnings compared to more 
serious warnings.In a qualified system, most warnings should exist in this form. 
The early warning system only needs to be reasonable because it is only an auxiliary 
means. Modern software often embeds this part and tries to expose it. This is 
actually a means of program optimization. Different people often use it to bring 
different results, which helps maintain and improve it.Reasonable program configuration is obviously very important. In order to enable the program to run in different environments, 
I try to configure them reasonably. There is no limit to the configuration file. In fact,
any configuration file with a suffix will work when there is a corresponding interpreter. 
Most of the time, there is no reason to choose to write an unnecessary interpreter just 
to interpret a configuration file. The general format of the dotenv file type can be used. 
And load these variables about program configuration in the most concise way to achieve 
multi-mode.
```