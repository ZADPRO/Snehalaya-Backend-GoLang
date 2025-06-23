package query

var AdminLoginSQL = `
SELECT
  u."refUserId",
  u."refUserCustId",
  u."refRTId",
  u."refUserFName",
  u."refUserLName",
  u."refUserBranchId",
  uac."refUACPassword",
  uac."refUACHashedPassword",
  uac."refUACUsername",
  ucd."refUCDMobile",
  ucd."refUCDEmail"
  -- *
FROM
  public."Users" u
  JOIN public."refUserAuthCred" uac ON u."refUserId" = uac."refUserId"
  JOIN public."refUserCommunicationDetails" ucd ON u."refUserId" = ucd."refUserId"
WHERE
  ucd."refUCDMobile" = $1
  AND u."refUserStatus" = 'true';
`
