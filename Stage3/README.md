# GoLangLab stage 3
### Extra Json in GO

You will adopt and continue working on your application constructed in Stage 1. You need to hand in the
complete and extended application including the following:



A function that stores the data sets from the JSON source in tables in a PostgreSQL database.
##### Status: Complete
A function that provides the same overview as requested in Stage 2, but you will now realise this
function using SQL queries to the PostgreSQL database
##### Status: Complete
Proof that the same output is produced in Stage 2 and Stage 3.
![Proof](https://cdn.discordapp.com/attachments/615404143200305167/770020722260181022/unknown.png)
- Formatting is slightly different but its only because im rounding slighting differently
##### Status: Complete
Measurements that show the time required to produce the requested overview. Explain the order of complexity of your solution.
- Measurements are built into the program. The time taken function can be added anywhere that you want to see how long a function takes.
- I'm not super confident at Big O but I think we're probably looking at O(n) as each function is still interacting with each object usually only once.
##### Status: Complete
An in-depth explanation on the differences in timespans you have found between stages 2 and 3.
* This part should be easy to explain. The difference in time is likely due to the speed differences between pulling from a database and pulling from RAM. RAM being the much much faster option as this is literally what it's designed for. 
##### Status: Complete

