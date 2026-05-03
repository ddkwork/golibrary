package driver

import "golang.org/x/sys/windows"

type Privilege struct {
	adjusted map[string]bool
}

func NewPrivilege() *Privilege {
	return &Privilege{adjusted: make(map[string]bool)}
}

func (p *Privilege) Enable(name string) bool {
	if p.adjusted[name] {
		return true
	}
	var token windows.Token
	currentProcess, _ := windows.GetCurrentProcess()
	windows.OpenProcessToken(currentProcess, windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY, &token)
	defer token.Close()

	var luid windows.LUID
	windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr(name), &luid)
	windows.AdjustTokenPrivileges(token, false, &windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges: [1]windows.LUIDAndAttributes{{Luid: luid, Attributes: windows.SE_PRIVILEGE_ENABLED}},
	}, 0, nil, nil)
	p.adjusted[name] = true
	return true
}

func (p *Privilege) SeCreateToken() bool              { return p.Enable("SeCreateTokenPrivilege") }
func (p *Privilege) SeAssignPrimaryToken() bool        { return p.Enable("SeAssignPrimaryTokenPrivilege") }
func (p *Privilege) SeLockMemory() bool                { return p.Enable("SeLockMemoryPrivilege") }
func (p *Privilege) SeIncreaseQuota() bool             { return p.Enable("SeIncreaseQuotaPrivilege") }
func (p *Privilege) SeUnsolicitedInput() bool          { return p.Enable("SeUnsolicitedInputPrivilege") }
func (p *Privilege) SeMachineAccount() bool            { return p.Enable("SeMachineAccountPrivilege") }
func (p *Privilege) SeTcb() bool                       { return p.Enable("SeTcbPrivilege") }
func (p *Privilege) SeSecurity() bool                  { return p.Enable("SeSecurityPrivilege") }
func (p *Privilege) SeTakeOwnership() bool             { return p.Enable("SeTakeOwnershipPrivilege") }
func (p *Privilege) SeLoadDriver() bool                { return p.Enable("SeLoadDriverPrivilege") }
func (p *Privilege) SeSystemProfile() bool             { return p.Enable("SeSystemProfilePrivilege") }
func (p *Privilege) SeSystemtime() bool                { return p.Enable("SeSystemtimePrivilege") }
func (p *Privilege) SeProfileSingleProcess() bool      { return p.Enable("SeProfileSingleProcessPrivilege") }
func (p *Privilege) SeIncreaseBasePriority() bool      { return p.Enable("SeIncreaseBasePriorityPrivilege") }
func (p *Privilege) SeCreatePagefile() bool            { return p.Enable("SeCreatePagefilePrivilege") }
func (p *Privilege) SeCreatePermanent() bool           { return p.Enable("SeCreatePermanentPrivilege") }
func (p *Privilege) SeBackup() bool                    { return p.Enable("SeBackupPrivilege") }
func (p *Privilege) SeRestore() bool                   { return p.Enable("SeRestorePrivilege") }
func (p *Privilege) SeShutdown() bool                  { return p.Enable("SeShutdownPrivilege") }
func (p *Privilege) SeDebug() bool                     { return p.Enable("SeDebugPrivilege") }
func (p *Privilege) SeAudit() bool                     { return p.Enable("SeAuditPrivilege") }
func (p *Privilege) SeSystemEnvironment() bool         { return p.Enable("SeSystemEnvironmentPrivilege") }
func (p *Privilege) SeChangeNotify() bool              { return p.Enable("SeChangeNotifyPrivilege") }
func (p *Privilege) SeRemoteShutdown() bool            { return p.Enable("SeRemoteShutdownPrivilege") }
func (p *Privilege) SeUndock() bool                    { return p.Enable("SeUndockPrivilege") }
func (p *Privilege) SeSyncAgent() bool                 { return p.Enable("SeSyncAgentPrivilege") }
func (p *Privilege) SeEnableDelegation() bool          { return p.Enable("SeEnableDelegationPrivilege") }
func (p *Privilege) SeManageVolume() bool              { return p.Enable("SeManageVolumePrivilege") }
func (p *Privilege) SeImpersonate() bool               { return p.Enable("SeImpersonatePrivilege") }
func (p *Privilege) SeCreateGlobal() bool              { return p.Enable("SeCreateGlobalPrivilege") }

var privilege = NewPrivilege()
