# Columba
Columba is a self-hosted package tracking system that was designed and tested to run on a Raspberry Pi. The name comes from the suborder of birds that contains pigeons and doves, most notably, the carrier pigeon.

## Table of Contents

- [Features](#features)
- [Supported Trackers](#supported-trackers)
- [Installation](#installation)
- [Configuration Settings](#configuration-settings)
- [Updating Columba](#updating-columba)
- [Uninstalling Columba](#uninstalling-columba)
- [Roadmap](#roadmap)
- [Roadmap](#licenses-and-attribution)

## Features

![Columba Web App](docs/img/img-cardlayout.jpg)

*Columba* has a simple, mobile-friendly web interface for quickly viewing all your active packages. 

![Columba Web App](docs/img/img-events.jpg)

Tracking services are regularly polled in the background, even if no web pages are open, and new tracking events are recorded for the next time you pull up the website.

![Columba Dark Theme >](docs/img/img-darktheme.jpg)
There is a selectable dark or light mode setting, and 3 different UI layouts for different users' preferences. 

Adding a package requires only the tracking number and shipping service provider, but an optional label can be provided if desired. 

![Columba Dark Theme >](docs/img/img-tablelayout.jpg)
Packages that have been delivered are highlighted in green.

![Columba Dark Theme >](docs/img/img-delete.jpg)
Packages can be removed from the view to clean up the list of packages once they are no longer needed.

## Supported Trackers
*Columba* currently supports the following tracking services:

- UPS
- USPS

See the [Roadmap](#roadmap) section for planned future support.

## Installation

*Columba* was designed to be easy to install on a Raspberry Pi.The easiest way to install *Columba* is to use `wget` or `curl` to pull the `install.sh` script file from the GitHub repository and execute it.

**cURL Example**
```
curl -s https://raw.githubusercontent.com/littlehawk93/columba/main/install/install.sh | /bin/bash
```
**wget Example**
```
wget -q https://raw.githubusercontent.com/littlehawk93/columba/main/install/install.sh | /bin/bash
```

The install script will pull the latest version of *Columba* from the GitHub releases and create the necessary service account, directories, and files for Columba to run as a service on your Pi. Once the install script finishes, you can enable the service with the command:
```
systemctl enable columba
```
and start the webapp using:
```
service columba start
```

## Configuration Settings

Configuration settings are saved in YAML format in the file: `/etc/columba/config.yaml`. A typical installation will produce a configuration file that looks like this:
```
database:
    database: /var/data/columba/db.sqlite
web:
    bind:
    port: 80
webroot: /var/www/columba
minrefreshtime: 1800
bgupdatetime: 3600
```

Below are the configuration settings and their purposes

| Parent   | Property Name  | Description                                                                                                                                                                                                                                                                                        | Default |
|----------|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| database | database       | The path to the SQLite database file used to store application data. If the file does not exists, *Columba* will create a new empty database file when it launches                                                                                                                                 |         |
| web      | bind           | The IP address to bind the *Columba* web process to. Leave blank to bind to all IP addresses                                                                                                                                                                                                       |         |
| web      | port           | The TCP port to bind the *Columba* web process to                                                                                                                                                                                                                                                  |    80   |
|          | webroot        | The directory containing the web application static files                                                                                                                                                                                                                                          |         |
|          | minrefreshtime | The minimum time, in seconds, before *Columba* will attempt to scrape a tracking service for shipment updates. This is a safety tool to prevent tracking services from blocking your IP address due to spamming                                                                                    |         |
|          | bgupdatetime   | Time in seconds between background updates. Each time a background update triggers, *Columba* will scrape all active packages for updates and save any new tracking events, even if no clients are connected to the website. This setting ignores the minrefreshtime parameter, so use responsibly |         |

## Updating Columba

During installation, *Columba* saves a copy of the install script as the file `/etc/columba/update.sh`. If you run the script again, it will check the current installed version of Columba and update executable and web files if a newer version is found on GitHub. It will **not** overwrite your configuration settings or SQLite database files. 

## Uninstalling Columba

To uninstall *Columba*, you can execute the script file located at `/etc/columba/uninstall.sh`. This will remove all the default directories, web files, and executable and remove the service account from your Pi. If you have changed the SQLite database file location or web files locations from their default locations, the uninstall script will not remove those files.

## Roadmap

Support for other tracking services is planned, but ironically there is no estimated time frame for when new features will arrive. Currently, the next providers planned are:
- FedEx 
- AliExpress Shipping

If you have ideas or requests for other shipment trackers. Please feel free to make an issue on the github project page. 

## Licenses and Attribution

The Columba project wouldn't be possible without the following projects:

- [google/uuid](https://github.com/google/uuid) by [Google](https://github.com/google)

- [glebarez/go-sqlite](https://github.com/glebarez/go-sqlite) by [glebarez](https://github.com/glebarez)

- [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery) by [PuerkitoBio](https://github.com/PuerkitoBio)

- [gorilla/mux](https://github.com/gorilla/mux) by [Gorilla web toolkit](https://github.com/gorilla)

- [https://github.com/spf13/cobra](spf13/cobra) by [Steve Francia](https://github.com/spf13)

- [https://github.com/spf13/viper](spf13/viper) by [Steve Francia](https://github.com/spf13)

[PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery), [gorilla/mux](https://github.com/gorilla/mux), [google/uuid](https://github.com/google/uuid), and [glebarez/go-sqlite](https://github.com/glebarez/go-sqlite) are all available under the BSD-3 License:

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

* Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
* Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----

The [spf13/cobra](https://github.com/spf13/cobra) library code is available under the Apache 2.0 License:

Apache License

Version 2.0, January 2004

[http://www.apache.org/licenses/](http://www.apache.org/licenses/)

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.