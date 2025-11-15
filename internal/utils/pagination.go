package utils

// Paginate converts pageNumber+pageSize into IMAP-style 1-based sequence numbers.
// It ensures the range is always valid, even on the last page.
func Paginate(total, pageSize, pageNumber int) (from, to int) {

	// Default safety: never return 0,0
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	// Calculate "to" position (newest-first logic)
	to = total - pageSize*(pageNumber-1)

	// Calculate "from" position
	from = to - pageSize + 1

	// Clamp the values into valid IMAP ranges
	if to > total {
		to = total
	}
	if from < 1 {
		from = 1
	}
	if from > to {
		from = to
	}

	return from, to
}




/*

intution behind the pagination
emails 
1 - oldest email
2 - older
N - newest email

total = N = mbox.Messages // we want pages of emails, newest first

total = total messages in the mailbox 
pageSize = number of emails per page
pageNumber = which page we want starting from 1 = newest page

in db and many applications computer only get from and to sequence number for the page

the newest email always at page 1 and at position N and should come at top
so we need from and to 

from = to - pageSize + 1 (let pageNumb) // +1 to make the imp range inclusive
if pagesize = 1,then page 1 should start 10 emails before the newest

for page 2, 
we need to change to = previos from (that was in page 1) -1 
and from = to - pageSize + 1

and generalizing for any page Number
to = total - pageSize * (pageNumber - 1)

emails to skip = pageSize * (number of previous pages)
because if we draw block the to is just the emails to skip the pageNumber * pageSize gives total emails in till the end ofo the page for example if pageNumber is 2 and pageSize is 10 and total emails= 50 then pageSize * pageNumber = 20 which tells that from N there are 20 emails 


I see data like this as a sequence like 1 to n and then for pageSize for example 10 I draw boxes aaround 10 10 document or email whatever u wanna say each box is labelled 1 they are stack one over another like the oldest on top latest at the bottom

to calculate from and to,

first of all understand pageNumber * pagSize it gives total number of documents from that pageNumber till the bottom using this as reference

we can calculate from for any pageNumber that will be 

from = total - pageNumber * pageSize + 1

and for to we just need the total number of ducmebnt just before the pageNumber we are considering thats why -1 
*/ 