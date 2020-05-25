# QuoteBot

### What is QuoteBot ?

QuoteBot is a Discord bot written in Go, and is mainly a Golang learning project for its developers.
It is builded upon the [discordgo library](https://github.com/bwmarrin/discordgo) in Go 1.14.

QuoteBot allows you to store & have access to a book of user quotes from server members.
Each member has a dedicated quote book containing things that have been deemed funny / interesting by other members.


### How does it work ?

To save a message as a quote, one simply has to **react to the message with the 💾 emoji** ( :floppy_disk: )
Quotes are only available on the server they were taken from : you can only get a quote back on the server it was written.