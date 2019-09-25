package configs

// Code generated DO NOT EDIT.

// Policies auto-generated from asset file
const Policies string = `version: v1.beta1
policies:
  all:
    filter: "(objectClass=*)"
    attributes:
      - name: dn
      - name: objectClass
  users:
    filter: "(objectClass=inetOrgPerson)"
    attributes:
      - name: cn
      - name: sn
      - name: userPassword
  inetOrgPerson:
    filter: "(objectClass=inetOrgPerson)"
    attributes:
      - name: dn
      - name: cn
      - name: sn
#     - name: audio
#     - name: businessCategory
#     - name: carLicense
#     - name: departmentNumber
      - name: displayName
      - name: employeeNumber
      - name: employeeType
      - name: givenName
#     - name: homePhone
#     - name: homePostalAddress
#     - name: initials
#     - name: jpegPhoto
#     - name: labeledURI
      - name: mail
#     - name: manager
#     - name: mobile
#     - name: o
#     - name: pager
#     - name: photo
#     - name: roomNumber
#     - name: secretary
      - name: uid
#     - name: userCertificate
#     - name: x500uniqueIdentifier
#     - name: preferredLanguage
#     - name: userSMIMECertificate
#     - name: userPKCS12
# organizationalPerson
#     - name: title
#     - name: x121Address
#     - name: registeredAddress
#     - name: destinationIndicator
#     - name: preferredDeliveryMethod
#     - name: telexNumber
#     - name: teletexTerminalIdentifier
#     - name: internationaliSDNNumber
#     - name: facsimileTelephoneNumber
#     - name: street
#     - name: postOfficeBox
#     - name: postalCode
#     - name: postalAddress
#     - name: physicalDeliveryOfficeName
#     - name: ou
#     - name: st
#     - name: l
# person
#     - name: userPassword
#     - name: telephoneNumber
#     - name: seeAlso
#     - name: description
  organizationalPerson:
    filter: "(objectClass=organizationalPerson)"
    attributes:
      - name: dn
      - name: cn
      - name: sn
#     - name: title
#     - name: x121Address
#     - name: registeredAddress
#     - name: destinationIndicator
#     - name: preferredDeliveryMethod
#     - name: telexNumber
#     - name: teletexTerminalIdentifier
#     - name: internationaliSDNNumber
#     - name: facsimileTelephoneNumber
      - name: street
#     - name: postOfficeBox
      - name: postalCode
      - name: postalAddress
#     - name: physicalDeliveryOfficeName
#     - name: ou
      - name: st
      - name: l
# person
#     - name: userPassword
#     - name: telephoneNumber
#     - name: seeAlso
#     - name: description
  residentialPerson:
    filter: "(objectClass=residentialPerson)"
    attributes:
      - name: dn
      - name: cn
      - name: sn
      - name: l
#     - name: businessCategory
#     - name: x121Address
#     - name: registeredAddress
#     - name: destinationIndicator
#     - name: telexNumber
#     - name: teletexTerminalIdentifier
#     - name: internationaliSDNNumber
#     - name: facsimileTelephoneNumber
#     - name: preferredDeliveryMethod
      - name: street
#     - name: postOfficeBox
      - name: postalCode
      - name: postalAddress
#     - name: physicalDeliveryOfficeName
      - name: st
# person
#     - name: userPassword
#     - name: telephoneNumber
#     - name: seeAlso
#     - name: description
  person:
    filter: "(objectClass=person)"
    attributes:
      - name: dn
      - name: cn
      - name: sn
      - name: userPassword
      - name: telephoneNumber
#     - name: seeAlso
      - name: description
  account:
    filter: "(objectClass=account)"
    attributes:
      - name: dn
      - name: uid
      - name: description
#     - name: seeAlso
#     - name: localityName
#     - name: organizationName
#     - name: organizationalUnitName
#     - name: host
  posixAccount:
    filter: "(objectClass=posixAccount)"
    attributes:
      - name: dn
      - name: cn
      - name: uid
      - name: uidNumber
      - name: gidNumber
      - name: homeDirectory
#     - name: userPassword
#     - name: loginShell
#     - name: gecos
#     - name: description
  simpleSecurityObject:
    filter: "(objectClass=simpleSecurityObject)"
    attributes:
      - name: dn
      - name: userPassword
  organizationalRole:
    filter: "(objectClass=organizationalRole)"
    attributes:
      - name: dn
      - name: cn
#     - name: x121Address
#     - name: registeredAddress
#     - name: destinationIndicator
#     - name: telexNumber
#     - name: teletexTerminalIdentifier
#     - name: telephoneNumber
#     - name: internationaliSDNNumber
#     - name: facsimileTelephoneNumber
#     - name: seeAlso
#     - name: roleOccupant
#     - name: preferredDeliveryMethod
#     - name: street
#     - name: postOfficeBox
#     - name: postalCode
#     - name: postalAddress
#     - name: physicalDeliveryOfficeName
#     - name: ou
#     - name: st
#     - name: l
#     - name: description
  organizationalUnit:
    filter: "(objectClass=organizationalUnit)"
    attributes:
      - name: dn
      - name: ou
#     - name: userPassword
#     - name: searchGuide
#     - name: seeAlso
#     - name: businessCategory
#     - name: x121Address
#     - name: registeredAddress
#     - name: destinationIndicator
#     - name: preferredDeliveryMethod
#     - name: telexNumber
#     - name: teletexTerminalIdentifier
#     - name: telephoneNumber
#     - name: internationaliSDNNumber
#     - name: facsimileTelephoneNumber
#     - name: street
#     - name: postOfficeBox
#     - name: postalCode
#     - name: postalAddress
#     - name: physicalDeliveryOfficeName
#     - name: st
#     - name: l
#     - name: description
  organization:
    filter: "(objectClass=organization)"
    attributes:
      - name: dn
      - name: o
#     - name: userPassword
#     - name: searchGuide
#     - name: seeAlso
#     - name: businessCategory
#     - name: x121Address
#     - name: registeredAddress
#     - name: destinationIndicator
#     - name: preferredDeliveryMethod
#     - name: telexNumber
#     - name: teletexTerminalIdentifier
#     - name: telephoneNumber
#     - name: internationaliSDNNumber
#     - name: facsimileTelephoneNumber
#     - name: street
#     - name: postOfficeBox
#     - name: postalCode
#     - name: postalAddress
#     - name: physicalDeliveryOfficeName
#     - name: st
#     - name: l
#     - name: description
  dcObject:
    filter: "(objectClass=dcObject)"
    attributes:
      - name: dn
      - name: dc
  groupOfNames:
    filter: "(objectClass=groupOfNames)"
    attributes:
      - name: dn
      - name: cn
      - name: member
#     - name: businessCategory
#     - name: seeAlso
      - name: owner
#     - name: ou
#     - name: o
#     - name: description
  groupOfUniqueNames:
    filter: "(objectClass=groupOfUniqueNames)"
    attributes:
      - name: dn
      - name: cn
      - name: uniqueMember
#     - name: businessCategory
#     - name: seeAlso
      - name: owner
#     - name: ou
#     - name: o
#     - name: description
  country:
    filter: "(objectClass=country)"
    attributes:
      - name: dn
      - name: c
#     - name: searchGuide
#     - name: description
  friendlyCountry:
    filter: "(objectClass=friendlyCountry)"
    attributes:
      - name: dn
      - name: c
      - name: friendlyCountryName
# country
#     - name: searchGuide
#     - name: description
  device:
    filter: "(objectClass=device)"
    attributes:
      - name: dn
      - name: cn
#     - name: serialNumber
#     - name: seeAlso
#     - name: owner
#     - name: ou
#     - name: o
#     - name: l
#     - name: description

# applicationProcess
# applicationEntity
# locality
# dmd
# domain
# ...
`
