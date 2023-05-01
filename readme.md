# T.E.S.

### Theion Eidolon Socors 

At its core, the T.E.S. is a simple, yet powerful, tool which helps users understand (nearly) any resource.

## Overview
---
Integral components of T.E.S. are:
- A server, which hosts an LLM, and generates responses as a simple API.
- A client, qhich sends a query.
- A query, which consists of URL's and/or Files, and a question asked by a human.
- A resonse, which is a helpful answer from a personable assistant acting as an expert in whatever field the user's question pertained to.
---

## IMPORTANT!
Currently, there are no security solutions in place for keeping data private (especially if multiple decentralized servers are implemented). 
As such it would be wise NOT to transfer or disclose private information / sensitive data within the T.E.S. system.

## Accessibility / Ease of use
This project is intended to further the total knowledge of humanity, and as such needs to be as accessible and user-friendly as possible.
---
A means to this end, the T.E.S. ecosystem (T.E.S.e.):
- 📱 IOS client, which along with accepting the normal text-based input, also allows you to record audio.
- 📱 Android client, which along with accepting the normal text-based input, also allows you to record audio.
- 💻 Web client. 
- 💬 SMS-based interaction.
- 📧 Email-based interaction.

### More Servers?
As a solo developer I am burdened with limited server availability. As such, at some point in the future it might become feasible to implement a decentralized multi-server system wherein anybody can host an API instance, and clients connect to whichever server has the least load, or response time, or optimize certain model parameters, etc.


# Server structure
Each server has a queue which contains all unanswered queries, aka pending jobs. These queries contain Human input questions, along with files, URLs, and return destination information.
When the server has adequate resources, it grabs a job from the queue and generates a helpful response based on the given resources and human input. It then sends the resonse to the client, and repeats the process.

## Server addons (TODO)
- [ ] Allow for larger/longer textual resource input by addinga  preprocesssing step wherein the server takes manageable but large chunks of resources, compresses that down into its gist, and then combines all of the summarized chunks to form a response with the human input.
- [ ] Allow the server to generate wikipedia queries and add data from wikipedia to its resources when generating a helpful response.
