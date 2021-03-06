# GoLangLab stage 1
## Status: Complete
### Import and look-up information from a json source

You will have access to a collection of data in JSON format that contains information about companies and
students doing their internship at these companies. You need to hand in the following:

- An application that will import these data sets into appropriate Go data structures. You will include a function that produces an overview in text format, each line containing the following information: the student's first name, surname, company name, and contact email address for that company
##### Status: Complete
- Proof that all students have been allocated to a company. In case not all students are allocated, explain why and count the number of unallocated students. This needs to be an automated proof that can be applied to a different JSON file (with the same structure).
.- *I check to see if the company id on the students file links to a company by checking if the id is greater than 0 but less than or equal to the company array length*.
.-  *I also list all the interns that don't have companies and all the companies that don't have interns*
##### Status: Complete
- Measurements that show the time required to produce the requested overview. Explain how the time relates to the size of the data set. E.g. what happens if the data set is twice its current size?
.- *I have created a function that measures the time funtions take using the defer method. In this instance the data is loaded in linear time so doubling the size of the data set will double the time taken*
##### Status: Complete