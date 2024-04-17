# godap
-------

A minimalistic LDAP server in Go

Forked from [bradleypeabody/godap](https://github.com/bradleypeabody/godap)
and ported to a fully working Bazel usable package

## What is this thing?

This is the beginnings of LDAP serving functionality written in Go.

It is not intended to be or become a full LDAP server, it's just minimal bind
and simplistic search functionality.
It aims to be enough to implement authentication-related LDAP services for use
with Apache httpd or anything else that supports LDAP auth.

Theoretically more could be built out to approach the functionality of
a full LDAP server. I don't have time for that stuff.

## How to use it

The LDAP server can be simply built with Bazel

```bash
bazel build :ldap-mock
```

The artifact is located in `./bazel-bin/ldap-mock_/` and named `ldap-mock`.

To start the server with a custom address (defaults to `127.0.0.1:389`) use
the `-addr=127.0.0.99:389` command line flag.

```bash
./bazel-bin/ldap-mock_/ldap-mock -addr=127.0.0.99:389
# or use the bazel run command
bazel run :ldap-mock -- -addr=127.0.0.99:389
```

Finally use [`ldapsearch`](https://docs.ldap.com/ldap-sdk/docs/tool-usages/ldapsearch.html)
or any other tool to contact the LDAP server

```bash
ldapsearch -x \
-b "dc=example,dc=com" \
-H "ldap://127.0.0.99:389" \
-D "cn=user,dc=example,dc=com" \
-w "1234" \
'(uid=carl)'
```

As this is a very basic LDAP server, it simply return what it has been ask for.
With the ldapsearch command from above the infos of `carl` will be retured.
Using a different filter query than `(uid=<name>)` will result in a
`32 No such object` error.

```
# extended LDIF
#
# LDAPv3
# base <dc=example,dc=com> with scope subtree
# filter: (uid=carl)
# requesting: ALL
#

# carl, example.com
dn: cn=carl,dc=example,dc=com
sn: carl
cn: carl
uid: carl
homeDirectory: /home/carl
objectClass: top
objectClass: posixAccount
objectClass: inetOrgPerson

# search result
search: 2
result: 0 Success

# numResponses: 2
# numEntries: 1
```

The unit test can be executed with

```bash
bazel test //godap:go_default_test
```

## Why was this made?

Because the road to hell is paved with good intentions.

The short version of the story goes like this:
I hate LDAP. I used to love it. But I loved it for all the wrong reasons.
LDAP is supported as an authentication solution by many different pieces of
software. Aside from its de jure standard status, its wide deployment
cements it as a de facto standard as well.

However, just because it is a standard doesn't mean it is a great idea.

I'll admit that given its age LDAP has had a good run. I'm sure its
authors carefully considered how to construct the protocol and chose
ASN.1 and its encoding with all of wellest of well meaning intentions.

The trouble is that with today's Internet, LDAP is just a pain in the ass.
You can't call it from your browser. It's not human readable or easy
to debug. Tooling is often arcane and confusing. It's way more complicated
than what is needed for most simple authentication-only uses. (Yes, I know
there are many other uses than authentication - but it's often too complicated
for those too.)

Likely owing to the complexity of the protocol, there seems to be virtually
no easy to use library to implement the server side of the LDAP protocol
that isn't tied in with some complete directory server system; and certainly
not in a language as easy to "make it work" as Go.

So this means that if you are a web developer and you have a database table
with usernames and (hopefully properly salted and hashed) passwords in it, and you
have a third party application that supports LDAP authentication, you
can't easily use your own user data source and just make it "speak LDAP"
to this other application.

Well, this project provides the basic guts to make that work.
<a href="https://github.com/bradleypeabody/godap/blob/master/godap_test.go">Have a look at the test file</a> for an example of what to do.

In a way, with this project I embrace LDAP. In the sense that a wrestler
embraces the neck of his opponent to form a headlock, and the unruly
brute is muscled into submission.

Dependencies
------------

ASN.1+BER encoding and decoding is done with [asn1-ber](https://github.com/go-asn1-ber/asn1-ber)

It works, it's cool.

Disclaimers
-----------

1. This thing is still fairly rough. Haven't had time to polish it. But I
wanted to get it up on Github in case anyone else was interested. I've
been searching for something like this off and on for a while, finally
got around to writing it.

2. It's not impossible that in some places I'm violating pieces of the
LDAP spec. My goal with this project was to get it so LDAP can be
useful as an authentication mechanism on top of other
non-directory data sources. It does that, I'm happy.
Pull requests welcome to fix things that aren't perfect. License
is MIT so feel free to do whatever you want with the code.
