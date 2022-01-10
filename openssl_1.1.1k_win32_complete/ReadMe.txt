OpenSSL 1.1.1k 25 Mar 2021 (32-bit)
Copyright (c) 1998-2021 The OpenSSL Project. All rights reserved.

This installation package was created by Catalyst Development for use with
our SocketTools suite of Internet components for Windows. For more information
about SocketTools, visit https://sockettools.com

This software is provided free of charge and you do not have to own a license
for SocketTools to use OpenSSL toolkit. For more information about the
OpenSSL license, review the LICENSE.TXT file included with this distribution.

Although SocketTools does not use the OpenSSL libraries themselves, the toolkit
can be useful for generating keys, certificate signining requests (CSRs) and
creating certificates. It can also be useful as an external tool to get
information about secure connections.

Please note that we cannot provide technical support for the use of OpenSSL.
If you want to contribute to the OpenSSL project, report a security bug or
review open issues, visit https://www.openssl.org/community/

This is a compiled redistribution of the OpenSSL toolkit for the 32-bit Windows
platform. It is a statically linked build using the standard version 1.1.1g
source from https://github.com/openssl/openssl

This is cryptography software and as such, its use may be restricted depending
on any applicable laws in your country that govern encryption. You alone are
responsible for knowing your legal rights and obligations. If you have any
legal questions or concerns about using OpenSSL, please consult a lawyer.

During the installation process, you can specify where you want to install the
executable, libraries and support files. The default installation location on
Windows XP and other 32-bit Windows platforms is:

C:\Program Files\OpenSSL

When installed on 64-bit Windows, OpenSSL will run under WoW64 and the files
will be copied to the folder for 32-bit applications:

C:\Program Files (x86)\OpenSSL

It uses the default build options, with the exception that the OPENSSL_CONF
location is the common application folder, rather than the default common files
folder. This is where the configuration file openssl.cnf can be found.
On Windows XP this folder location is:

C:\Documents and Settings\All Users\Application Data\OpenSSL

On Windows 7, Windows 8 and Windows 10 this folder location is:

C:\ProgramData\OpenSSL

The OpenSSL documentation is formatted as HTML versions of UNIX "man" pages
and we have combined them into Microsoft Compiled HTML Help format which is
generally easier to use on Windows. Keep in mind that the documentation is
created during the build process and although functional, it is fairly basic
in its layout and presentation. In other words, don't expect it to look like
Visual Studio documentation. You can find the online version of the OpenSSL
documentation and FAQ at https://www.openssl.org/docs/

The documentation is only included in the complete installation package. The
minimal installation package only includes the OpenSSL tool itself which can
be used from the command line.
