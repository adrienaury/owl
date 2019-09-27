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
  # openldap core.schema (MAY attributes are commented)
  top:
    filter: "(objectClass=top)"
    attributes:
      - name: dn
      - name: objectClass
  alias:
    inherit: top
    filter: "(objectClass=alias)"
    attributes:
      - name: dn
      - name: aliasedObjectName
  country:
    inherit: top
    filter: "(objectClass=country)"
    attributes:
      - name: dn
      - name: c
#     - name: searchGuide
#     - name: description
  locality:
    inherit: top
    filter: "(objectClass=locality)"
    attributes:
      - name: dn
#     - name: street
#     - name: seeAlso
#     - name: searchGuide
#     - name: st
#     - name: l
#     - name: description
  organization:
    inherit: top
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
  organizationalUnit:
    inherit: top
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
  person:
    inherit: top
    filter: "(objectClass=person)"
    attributes:
      - name: dn
      - name: sn
      - name: cn
#     - name: userPassword
#     - name: telephoneNumber
#     - name: seeAlso
#     - name: description
  organizationalPerson:
    inherit: person
    filter: "(objectClass=organizationalPerson)"
    attributes:
      - name: dn
#     - name: title
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
#     - name: ou
#     - name: st
#     - name: l
  organizationalRole:
    inherit: top
    filter: "(objectClass=organizationalRole)"
    attributes:
      - name: dn
      - name: cn
#     - name: x121Address
#     - name: registeredAddress
#     - name: destinationIndicator
#     - name: preferredDeliveryMethod
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
  groupOfNames:
    inherit: top
    filter: "(objectClass=groupOfNames)"
    attributes:
      - name: dn
      - name: member
      - name: cn
#     - name: businessCategory
#     - name: seeAlso
#     - name: owner
#     - name: ou
#     - name: o
#     - name: description
  residentialPerson:
    inherit: person
    filter: "(objectClass=residentialPerson)"
    attributes:
      - name: dn
      - name: l
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
  applicationProcess:
    inherit: top
    filter: "(objectClass=applicationProcess)"
    attributes:
      - name: dn
      - name: cn
#     - name: seeAlso
#     - name: ou
#     - name: l
#     - name: description
  applicationEntity:
    inherit: top
    filter: "(objectClass=applicationEntity)"
    attributes:
      - name: dn
      - name: presentationAddress
      - name: cn
#     - name: supportedApplicationContext
#     - name: seeAlso
#     - name: ou
#     - name: l
#     - name: description
  dSA:
    inherit: applicationEntity
    filter: "(objectClass=dSA)"
    attributes:
      - name: dn
#     - name: knowledgeInformation
  device:
    inherit: top
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
  strongAuthenticationUser:
    inherit: top
    filter: "(objectClass=strongAuthenticationUser)"
    attributes:
      - name: dn
      - name: userCertificate
  certificationAuthority:
    inherit: top
    filter: "(objectClass=certificationAuthority)"
    attributes:
      - name: dn
      - name: authorityRevocationList
      - name: certificateRevocationList
      - name: cACertificate
#     - name: crossCertificatePair
  groupOfUniqueNames:
    inherit: top
    filter: "(objectClass=groupOfUniqueNames)"
    attributes:
      - name: dn
      - name: uniqueMember
      - name: cn
#     - name: businessCategory
#     - name: seeAlso
#     - name: owner
#     - name: ou
#     - name: o
#     - name: description
  userSecurityInformation:
    inherit: top
    filter: "(objectClass=userSecurityInformation)"
    attributes:
      - name: dn
#     - name: supportedAlgorithms
  certificationAuthority-V2:
    inherit: certificationAuthority
    filter: "(objectClass=certificationAuthority-V2)"
    attributes:
      - name: dn
#     - name: deltaRevocationList
  cRLDistributionPoint:
    inherit: top
    filter: "(objectClass=cRLDistributionPoint)"
    attributes:
      - name: dn
      - name: cn
#     - name: certificateRevocationList
#     - name: authorityRevocationList
#     - name: deltaRevocationList
  dmd:
    inherit: top
    filter: "(objectClass=dmd)"
    attributes:
      - name: dn
      - name: dmdName
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
  extensibleObject:
    inherit: top
    filter: "(objectClass=extensibleObject)"
    attributes:
      - name: dn
  labeledURIObject:
    inherit: top
    filter: "(objectClass=labeledURIObject)"
    attributes:
      - name: dn
#     - name: labeledURI
  dynamicObject:
    inherit: top
    filter: "(objectClass=dynamicObject)"
    attributes:
      - name: dn
  simpleSecurityObject:
    inherit: top
    filter: "(objectClass=simpleSecurityObject)"
    attributes:
      - name: dn
      - name: userPassword
  dcObject:
    inherit: top
    filter: "(objectClass=dcObject)"
    attributes:
      - name: dn
      - name: dc
  uidObject:
    inherit: top
    filter: "(objectClass=uidObject)"
    attributes:
      - name: dn
      - name: uid
  referral:
    inherit: top
    filter: "(objectClass=referral)"
    attributes:
      - name: dn
      - name: ref
  # openldap cosine.schema (MAY attributes are commented)
  pilotPerson:
    inherit: top
    filter: "(objectClass=pilotPerson)"
    attributes:
      - name: dn
#     - name: userid
#     - name: textEncodedORAddress
#     - name: rfc822Mailbox
#     - name: favouriteDrink
#     - name: roomNumber
#     - name: userClass
#     - name: homeTelephoneNumber
#     - name: homePostalAddress
#     - name: secretary
#     - name: personalTitle
#     - name: preferredDeliveryMethod
#     - name: businessCategory
#     - name: janetMailbox
#     - name: otherMailbox
#     - name: mobileTelephoneNumber
#     - name: pagerTelephoneNumber
#     - name: organizationalStatus
#     - name: mailPreferenceOption
#     - name: personalSignature
  account:
    inherit: top
    filter: "(objectClass=account)"
    attributes:
      - name: dn
      - name: userid
#     - name: description
#     - name: seeAlso
#     - name: localityName
#     - name: organizationName
#     - name: organizationalUnitName
#     - name: host
  document:
    inherit: top
    filter: "(objectClass=document)"
    attributes:
      - name: dn
      - name: documentIdentifier
#     - name: commonName
#     - name: description
#     - name: seeAlso
#     - name: localityName
#     - name: organizationName
#     - name: organizationalUnitName
#     - name: documentTitle
#     - name: documentVersion
#     - name: documentAuthor
#     - name: documentLocation
#     - name: documentPublisher
  room:
    inherit: top
    filter: "(objectClass=room)"
    attributes:
      - name: dn
      - name: commonName
#     - name: roomNumber
#     - name: description
#     - name: seeAlso
#     - name: telephoneNumber
  documentSeries:
    inherit: top
    filter: "(objectClass=documentSeries)"
    attributes:
      - name: dn
      - name: commonName
#     - name: description
#     - name: seeAlso
#     - name: telephonenumber
#     - name: localityName
#     - name: organizationName
#     - name: organizationalUnitName
  domain:
    inherit: top
    filter: "(objectClass=domain)"
    attributes:
      - name: dn
      - name: domainComponent
#     - name: associatedName
#     - name: organizationName
#     - name: description
#     - name: businessCategory
#     - name: seeAlso
#     - name: searchGuide
#     - name: userPassword
#     - name: localityName
#     - name: stateOrProvinceName
#     - name: streetAddress
#     - name: physicalDeliveryOfficeName
#     - name: postalAddress
#     - name: postalCode
#     - name: postOfficeBox
#     - name: streetAddress
#     - name: facsimileTelephoneNumber
#     - name: internationalISDNNumber
#     - name: telephoneNumber
#     - name: teletexTerminalIdentifier
#     - name: telexNumber
#     - name: preferredDeliveryMethod
#     - name: destinationIndicator
#     - name: registeredAddress
#     - name: x121Address
  RFC822localPart:
    inherit: domain
    filter: "(objectClass=RFC822localPart)"
    attributes:
      - name: dn
#     - name: commonName
#     - name: surname
#     - name: description
#     - name: seeAlso
#     - name: telephonenumber
#     - name: physicalDeliveryOfficeName
#     - name: postalAddress
#     - name: postalCode
#     - name: postOfficeBox
#     - name: streetAddress
#     - name: facsimileTelephoneNumber
#     - name: internationalISDNNumber
#     - name: telephoneNumber
#     - name: teletexTerminalIdentifier
#     - name: telexNumber
#     - name: preferredDeliveryMethod
#     - name: destinationIndicator
#     - name: registeredAddress
#     - name: x121Address
  dNSDomain:
    inherit: domain
    filter: "(objectClass=dNSDomain)"
    attributes:
      - name: dn
#     - name: ARecord
#     - name: MDRecord
#     - name: MXRecord
#     - name: NSRecord
#     - name: SOARecord
#     - name: CNAMERecord
  domainRelatedObject:
    inherit: top
    filter: "(objectClass=domainRelatedObject)"
    attributes:
      - name: dn
      - name: associatedDomain
  friendlyCountry:
    inherit: country
    filter: "(objectClass=friendlyCountry)"
    attributes:
      - name: dn
      - name: friendlyCountryName
  pilotOrganization:
    inherit:
      - organization
      - organizationalUnit
    filter: "(objectClass=pilotOrganization)"
    attributes:
      - name: dn
#     - name: buildingName
  pilotDSA:
    inherit: dsa
    filter: "(objectClass=pilotDSA)"
    attributes:
      - name: dn
#     - name: dSAQuality
  qualityLabelledData:
    inherit: top
    filter: "(objectClass=qualityLabelledData)"
    attributes:
      - name: dn
      - name: dsaQuality
#     - name: subtreeMinimumQuality
#     - name: subtreeMaximumQuality
  # openldap cosine.schema (MAY attributes are commented)
  inetOrgPerson:
    inherit: organizationalPerson
    filter: "(objectClass=inetOrgPerson)"
    attributes:
      - name: dn
#     - name: audio
#     - name: businessCategory
#     - name: carLicense
#     - name: departmentNumber
#     - name: displayName
#     - name: employeeNumber
#     - name: employeeType
#     - name: givenName
#     - name: homePhone
#     - name: homePostalAddress
#     - name: initials
#     - name: jpegPhoto
#     - name: labeledURI
#     - name: mail
#     - name: manager
#     - name: mobile
#     - name: o
#     - name: pager
#     - name: photo
#     - name: roomNumber
#     - name: secretary
#     - name: uid
#     - name: userCertificate
#     - name: x500uniqueIdentifier
#     - name: preferredLanguage
#     - name: userSMIMECertificate
#     - name: userPKCS12
  # openldap nis.schema (MAY attributes are commented)
  posixAccount:
    inherit: top
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
  shadowAccount:
    inherit: top
    attributes:
      - name: dn
      - name: uid
#     - name: userPassword
#     - name: shadowLastChange
#     - name: shadowMin
#     - name: shadowMax
#     - name: shadowWarning
#     - name: shadowInactive
#     - name: shadowExpire
#     - name: shadowFlag
#     - name: description
  posixGroup:
    inherit: top
    attributes:
      - name: dn
      - name: cn
      - name: gidNumber
#     - name: userPassword
#     - name: memberUid
#     - name: description
  ipService:
    inherit: top
    attributes:
      - name: dn
      - name: cn
      - name: ipServicePort
      - name: ipServiceProtocol
#     - name: description
  ipProtocol:
    inherit: top
    attributes:
      - name: dn
      - name: cn
      - name: ipProtocolNumber
      - name: description
#     - name: description
  oncRpc:
    inherit: top
    attributes:
      - name: dn
      - name: cn
      - name: oncRpcNumber
      - name: description
#     - name: description
  ipHost:
    inherit: top
    attributes:
      - name: dn
      - name: cn
      - name: ipHostNumber
#     - name: l
#     - name: description
#     - name: manager
  ipNetwork:
    inherit: top
    attributes:
      - name: dn
      - name: cn
      - name: ipNetworkNumber
#     - name: ipNetmaskNumber
#     - name: l
#     - name: description
#     - name: manager
  nisNetgroup:
    inherit: top
    attributes:
      - name: dn
      - name: cn
#     - name: nisNetgroupTriple
#     - name: memberNisNetgroup
#     - name: description
  nisMap:
    inherit: top
    attributes:
      - name: dn
      - name: nisMapName
#     - name: description
  nisObject:
    inherit: top
    attributes:
      - name: dn
      - name: cn
      - name: nisMapEntry
      - name: nisMapName
#     - name: description
  ieee802Device:
    inherit: top
    attributes:
      - name: dn
#     - name: macAddress
  bootableDevice:
    inherit: top
    attributes:
      - name: dn
#     - name: bootFile
#     - name: bootParameter
`
