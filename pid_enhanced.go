package l5x

// PID_ENHANCED Reference: https://literature.rockwellautomation.com/idc/groups/literature/documents/rm/1756-rm006_-en-p~k.pdf
var pid_enhancedDT = DataType{
	Name: "PID_ENHANCED",
	Members: []Member{
		// Enable input.
		// If false, the instruction does not execute and outputs are not updated.
		// Default is true.
		{Name: "EnableIn", DataType: "BOOL"},

		// Scaled process variable input.
		// This value is typically read from an analog input module.
		// Valid = any float.
		// Default = 0.0
		{Name: "PV", DataType: "REAL"},

		// PV bad health indicator.
		// If PV is read from an analog input, then PVFault is normally controlled by the analog input fault status.
		// When PVFault is true, it indicates the input signal has an error.
		// Default is false = "good health"
		{Name: "PVFault", DataType: "BOOL"},

		// Maximum scaled value for PV.
		// The value of PV and SP which corresponds to 100 percent span of the Process Variable.
		// Valid = PVEUMin < PVEUMax <= maximum positive float.
		// Default = 100.0
		{Name: "PVEUMax", DataType: "REAL"},

		// Minimum scaled value for PV.
		// The value of PV and SP which corresponds to 0 percent span of the Process Variable.
		// Valid = maximum negative float <= PVEUMin < PVEUMax Default = 0.0
		{Name: "PVEUMin", DataType: "REAL"},

		// SP program value, scaled in PV units.
		// SP is set to this value when in Program control and not Cascade/Ratio mode.
		// If the value of SPProg < SPLLimit or > SPHLimit, the instruction sets the appropriate bit in Status and limits the value used for SP.
		// Valid = SPLLimit to SPHLimit.
		// Default = 0.0
		{Name: "SPProg", DataType: "REAL"},

		// SP operator value, scaled in PV units.
		// SP is set to this value when in Operator control and not Cascade/Ratio mode.
		// If the value of SPOper < SPLLimit or > SPHLimit, the instruction sets the appropriate bit in Status and limits the value used for SP.
		// Valid = SPLLimit to SPHLimit.
		// Default = 0.0
		{Name: "SPOper", DataType: "REAL"},

		// SP Cascade value, scaled in PV units.
		// If CascadeRatio is true and UseRatio is false, then SP = SPCascade.
		// This is typically the CVEU of a primary loop.
		// If CascadeRatio and UseRatio are true, then SP = (SPCascade x Ratio).
		// If the value of SPCascade < SPLLimit or > SPHLimit, set the appropriate bit in Status and limit the value used for SP.
		// Valid = SPLLimit to SPHLimit Default = 0.0
		{Name: "SPCascade", DataType: "REAL"},

		// SP high limit value, scaled in PV units.
		// If SPHLimit > PVEUMax, the instruction sets the appropriate bit in Status.
		// Valid = SPLLimit to PVEUMax.
		// Default = 100.0
		{Name: "SPHLimit", DataType: "REAL"},

		// SP low limit value, scaled in PV units.
		// If SPLLimit < PVEUMin, the instruction sets the appropriate bit in Status.
		// If SPHLimit < SPLLimit, the instruction sets the appropriate bit in Status and limits SP using the value of SPLLimit.
		// Valid = PVEUMin to SPHLimit.
		// Default = 0.0
		{Name: "SPLLimit", DataType: "REAL"},

		// Allow ratio control permissive.
		// Set to true to enable ratio control when in Cascade/Ratio mode.
		// Default is false.
		{Name: "UseRatio", DataType: "BOOL"},

		// Ratio program multiplier.
		// Ratio and RatioOper are set equal to this value when in Program control.
		// If RatioProg < RatioLLimit or > RatioHLimit, the instruction sets the appropriate bit in Status and limits the value used for Ratio.
		// Valid = RatioLLimit to RatioHLimit.
		// Default = 1.0
		{Name: "RatioProg", DataType: "REAL"},

		// Ratio operator multiplier.
		// Ratio is set equal to this value when in Operator control.
		// If RatioOper < RatioLLimit or > RatioHLimit, the instruction sets the appropriate bit in Status and limits the value used for Ratio.
		// Valid = RatioLLimit to RatioHLimit.
		// Default = 1.0
		{Name: "RatioOper", DataType: "REAL"},

		// Ratio high limit value.
		// Limits the value of Ratio obtained from RatioProg or RatioOper.
		// If RatioHLimit < RatioLLimit, the instruction sets the appropriate bit in Status and limits Ratio using the value of RatioLLimit.
		// Valid = RatioLLimit to maximum positive float.
		// Default = 1.0
		{Name: "RatioHLimit", DataType: "REAL"},

		// Ratio low limit value.
		// Limits the value of Ratio obtained from RatioProg or RatioOper.
		// If RatioLLimit < 0, the instruction sets the appropriate bit in Status and limits the value to zero.
		// If RatioHLimit < RatioLLimit, the instruction sets the appropriate bit in Status and limits Ratio using the value of RatioLLimit.
		// Valid = 0.0 to RatioHLimit.
		// Default = 1.0
		{Name: "RatioLLimit", DataType: "REAL"},

		// Control variable bad health indicator.
		// If CVEU controls an analog output, then CVFault normally comes from the analog output’s fault status.
		// When true, CVFault indicates an error on the output module and the instruction sets the appropriate bit in Status.
		// Default is false = "good health"
		{Name: "CVFault", DataType: "BOOL"},

		// CV initialization request.
		// This signal is normally controlled by the "In Hold" status on the analog output module controlled by CVEU or from the InitPrimary output of a secondary PID loop.
		// Default is false.
		{Name: "CVInitReq", DataType: "BOOL"},

		// CVEU initialization value, scaled in CVEU units.
		// When CVInitializing is true, CVEU = CVInitValue and CV equals the corresponding percentage value.
		// CVInitValue comes from the feedback of the analog output controlled by CVEU or from the setpoint of a secondary loop.
		// Instruction initialization is disabled when CVFaulted or CVEUSpanInv is true.
		// Valid = any float.
		// Default = 0.0
		{Name: "CVInitValue", DataType: "REAL"},

		// CV program manual value.
		// CV equals this value when in Program Manual mode.
		// If CVProg < 0 or > 100, or < CVLLimit or > CVHLimit when CVManLimiting is true, the instruction sets the appropriate bit in Status and limits the CV value.
		// Valid = 0.0 to 100.0.
		// Default = 0.0
		{Name: "CVProg", DataType: "REAL"},

		// CV operator manual value.
		// CV equals this value when in Operator Manual mode.
		// If not Operator Manual mode, the instruction sets CVOper = CV at the end of each instruction execution.
		// If CVOper < 0 or > 100, or < CVLLimit or > CVHLimit when CVManLimiting is true, the instruction sets the appropriate bit in Status and limits the CV value.
		// Valid = 0.0 to 100.0
		// Default = 0.0
		{Name: "CVOper", DataType: "REAL"},

		// CV override value.
		// CV equals this value when in override mode.
		// This value should correspond to a safe state output of the PID loop.
		// If CVOverride < 0 or >100, the instruction sets the appropriate bit in Status and limits the CV value.
		// Valid = 0.0 to 100.0
		// Default = 0.0
		{Name: "CVOverride", DataType: "REAL"},

		// CV_{n-1} value.
		// If CVSetPrevious is set, CV_{n-1} equals this value.
		// CV_{n-1} is the value of CV from the previous execution.
		// CVPrevious is ignored when in manual, override or hand mode or when CVInitializing is set.
		// If CVPrevious < 0 or > 100, or < CVLLimit or > CVHLimit when in Auto or cascade/ratio mode, the instruction sets the appropriate bit in Status and limits the CVn-1 value.
		// Valid = 0.0 to 100.0
		// Default = 0.0
		{Name: "CVPrevious", DataType: "REAL"},

		//Request to use CVPrevious. If true, CVn-1 = CVPrevious Default is false.
		{Name: "CVSetPrevious", DataType: "BOOL"},

		//Limit CV in manual mode request. If Manual mode and CVManLimiting is true, CV is limited by the CVHLimit and CVLLimit values.
		// Default is false.
		{Name: "CVManLimiting", DataType: "BOOL"},

		//Maximum value for CVEU. The value of CVEU which corresponds to 100 percent CV. If CVEUMax = CVEUMin, the instruction sets the appropriate bit in Status.
		// Valid = any float
		// Default = 100.0
		{Name: "CVEUMax", DataType: "REAL"},

		//Minimum value of CVEU. The value of CVEU which corresponds to 0 percent CV. If CVEUMax = CVEUMin, the instruction sets the appropriate bit in Status.
		// Valid = any float
		// Default = 0.0
		{Name: "CVEUMin", DataType: "REAL"},

		//CV high limit value. This is used to set the CVHAlarm output. It is also used for limiting CV when in Auto or Cascade/Ratio mode, or Manual mode if CVManLimiting is true. If CVHLimit > 100 or < CVLLimit, the instruction sets the appropriate bit in Status. If CVHLimit < CVLLimit, the instruction limits CV using the value of CVLLimit.
		// Valid = CVLLimit < CVHLimit <= 100.0
		// Default = 100.0
		{Name: "CVHLimit", DataType: "REAL"},

		//CV low limit value. This is used to set the CVLAlarm output. It is also used for limiting CV when in Auto or Cascade/Ratio mode, or Manual mode if CVManLimiting is true. If CVLLimit < 0 or CVHLimit < CVLLimit, the instruction sets the appropriate bit in Status. If CVHLimit < CVLLimit, the instruction limits CV using the value of CVLLimit.
		// Valid = 0.0 <= CVLLimit < CVHLimit
		// Default = 0.0
		{Name: "CVLLimit", DataType: "REAL"},

		//CV rate of change limit, in percent per second. Rate of change limiting is only used when in Auto or Cascade/Ratio modes or Manual mode if CVManLimiting is true. Enter 0 to disable CV ROC limiting. If CVROCLimit < 0, the instruction sets the appropriate bit in Status and disables CV ROC limiting.
		// Valid = 0.0 to maximum positive float
		// Default = 0.0
		{Name: "CVROCLimit", DataType: "REAL"},

		//Feed forward value. The value of feed forward is summed with CV after the zero- crossing deadband limiting has been applied to CV. Therefore changes in FF are always reflected in the final output value of CV. If FF < -100 or > 100, the instruction sets the appropriate bit in Status and limits the value used for FF.
		// Valid = -100.0 to 100.0
		// Default = 0.0
		{Name: "FF", DataType: "REAL"},

		//FF_{n-1} value. If FF SetPrevous is set, the instruction sets FF_{n-1} = FFPrevious. FF_{n-1} is the valu eof FF from the previous execution. If FFPrevious < -100 or > 100, the instruction sets the appropriate bit in Status and limits value used for FFn-1
		// Valid = -100.0 to 100.0
		// Default - 0.0
		{Name: "FFPrevious", DataType: "REAL"},

		//Request to use FFPrevious. If true, FFn-1 = FFPrevious. Default is false.
		{Name: "FFSetPrevious", DataType: "BOOL"},

		//CV Hand feedback value. CV equals this value when in Hand mode and HandFBFault is false (good health). This value typically comes from the output of a field mounted hand/ auto station and is used to generate a bumpless transfer out of hand mode. If HandFB < 0 or > 100, the instruction sets the appropriate bit in Status and limits the value used for CV.
		// Valid = 0.0 to 100.0 Default = 0.0
		{Name: "HandFB", DataType: "REAL"},

		//HandFB value bad health indicator. If the HandFB value is read from an analog input, then HandFBFault is typically controlled by the status of the analog input channel. When true, HandFBFault indicates an error on the input module and the instruction sets the appropriate bit in Status.
		// Default is false = "good health"
		{Name: "HandFBFault", DataType: "BOOL"},

		//Windup high request. When true, the CV cannot integrate in a positive direction. The signal is typically obtained from the WindupHOut output from a secondary loop. Default is false.
		{Name: "WindupHIn", DataType: "BOOL"},

		//Windup low request. When true, the CV cannot integrate in a negative direction. This signal is typically obtained from the WindupLOut output from a secondary loop. Default is false.
		{Name: "WindupLIn", DataType: "BOOL"},

		//
		{Name: "ControlAction", DataType: "BOOL"},

		//
		{Name: "DependIndepend", DataType: "BOOL"},

		//
		{Name: "PGain", DataType: "REAL"},

		//
		{Name: "IGain", DataType: "REAL"},

		//
		{Name: "DGain", DataType: "REAL"},

		//
		{Name: "PVEProportional", DataType: "BOOL"},

		//
		{Name: "PVEDerivative", DataType: "BOOL"},

		//
		{Name: "DSmoothing", DataType: "BOOL"},

		//
		{Name: "PVTracking", DataType: "BOOL"},

		//
		{Name: "ZCDeadband", DataType: "REAL"},

		//
		{Name: "ZCOff", DataType: "BOOL"},

		//
		{Name: "PVHHLimit", DataType: "REAL"},

		//
		{Name: "PVHLimit", DataType: "REAL"},

		//
		{Name: "PVLLimit", DataType: "REAL"},

		//
		{Name: "PVLLLimit", DataType: "REAL"},

		//
		{Name: "PVDeadband", DataType: "REAL"},

		//
		{Name: "PVROCPosLimit", DataType: "REAL"},

		//
		{Name: "PVROCNegLimit", DataType: "REAL"},

		//
		{Name: "PVROCPeriod", DataType: "REAL"},

		//
		{Name: "DevHHLimit", DataType: "REAL"},

		//
		{Name: "DevHLimit", DataType: "REAL"},

		//
		{Name: "DevLLimit", DataType: "REAL"},

		//
		{Name: "DevLLLimit", DataType: "REAL"},

		//
		{Name: "DevDeadband", DataType: "REAL"},

		//
		{Name: "AllowCasRat", DataType: "BOOL"},

		//
		{Name: "ManualAfterInit", DataType: "BOOL"},

		//
		{Name: "ProgProgReq", DataType: "BOOL"},

		//
		{Name: "ProgOperReq", DataType: "BOOL"},

		//
		{Name: "ProgCasRatReq", DataType: "BOOL"},

		//
		{Name: "ProgAutoReq", DataType: "BOOL"},

		//
		{Name: "ProgManualReq", DataType: "BOOL"},

		//
		{Name: "ProgOverrideReq", DataType: "BOOL"},

		//
		{Name: "ProgHandReq", DataType: "BOOL"},

		//
		{Name: "OperProgReq", DataType: "BOOL"},

		//
		{Name: "OperOperReq", DataType: "BOOL"},

		//
		{Name: "OperCasRatReq", DataType: "BOOL"},

		//
		{Name: "OperAutoReq", DataType: "BOOL"},

		//
		{Name: "OperManualReq", DataType: "BOOL"},

		//
		{Name: "ProgValueReset", DataType: "BOOL"},

		//
		{Name: "TimingMode", DataType: "DINT"},

		//
		{Name: "OversampleDT", DataType: "REAL"},

		//
		{Name: "RTSTime", DataType: "DINT"},

		//
		{Name: "RTSTimeStamp", DataType: "DINT"},

		//
		{Name: "AtuneAcquire", DataType: "BOOL"},

		//
		{Name: "AtuneStart", DataType: "BOOL"},

		//
		{Name: "AtuneUseGains", DataType: "BOOL"},

		//
		{Name: "AtuneAbort", DataType: "BOOL"},

		//
		{Name: "AtuneUnacquire", DataType: "BOOL"},

		//
		{Name: "EnableOut", DataType: "BOOL"},

		//
		{Name: "CVEU", DataType: "REAL"},

		//
		{Name: "CV", DataType: "REAL"},

		//
		{Name: "CVInitializing", DataType: "BOOL"},

		//
		{Name: "CVHAlarm", DataType: "BOOL"},

		//
		{Name: "CVLAlarm", DataType: "BOOL"},

		//
		{Name: "CVROCAlarm", DataType: "BOOL"},

		//
		{Name: "SP", DataType: "REAL"},

		//
		{Name: "SPPercent", DataType: "REAL"},

		//
		{Name: "SPHAlarm", DataType: "BOOL"},

		//
		{Name: "SPLAlarm", DataType: "BOOL"},

		//
		{Name: "PVPercent", DataType: "REAL"},

		//
		{Name: "E", DataType: "REAL"},

		//
		{Name: "EPercent", DataType: "REAL"},

		//
		{Name: "InitPrimary", DataType: "BOOL"},

		//
		{Name: "WindupHOut", DataType: "BOOL"},

		//
		{Name: "WindupLOut", DataType: "BOOL"},

		//
		{Name: "Ratio", DataType: "REAL"},

		//
		{Name: "RatioHAlarm", DataType: "BOOL"},

		//
		{Name: "RatioLAlarm", DataType: "BOOL"},

		//
		{Name: "ZCDeadbandOn", DataType: "BOOL"},

		//
		{Name: "PVHHAlarm", DataType: "BOOL"},

		//
		{Name: "PVHAlarm", DataType: "BOOL"},

		//
		{Name: "PVLAlarm", DataType: "BOOL"},

		//
		{Name: "PVLLAlarm", DataType: "BOOL"},

		//
		{Name: "PVROCPosAlarm", DataType: "BOOL"},

		//
		{Name: "PVROCNegAlarm", DataType: "BOOL"},

		//
		{Name: "DevHHAlarm", DataType: "BOOL"},

		//
		{Name: "DevHAlarm", DataType: "BOOL"},

		//
		{Name: "DevLAlarm", DataType: "BOOL"},

		//
		{Name: "DevLLAlarm", DataType: "BOOL"},

		//
		{Name: "ProgOper", DataType: "BOOL"},

		//
		{Name: "CasRat", DataType: "BOOL"},

		//
		{Name: "Auto", DataType: "BOOL"},

		//
		{Name: "Manual", DataType: "BOOL"},

		//
		{Name: "Override", DataType: "BOOL"},

		//
		{Name: "Hand", DataType: "BOOL"},

		//
		{Name: "DeltaT", DataType: "REAL"},

		//
		{Name: "AtuneReady", DataType: "BOOL"},

		//
		{Name: "AtuneOn", DataType: "BOOL"},

		//
		{Name: "AtuneDone", DataType: "BOOL"},

		//
		{Name: "AtuneAborted", DataType: "BOOL"},

		//
		{Name: "AtuneBusy", DataType: "BOOL"},

		//
		{Name: "Status1", DataType: "DINT"},

		//
		{Name: "Status2", DataType: "DINT"},

		//
		{Name: "InstructFault", DataType: "BOOL"},

		//
		{Name: "PVFaulted", DataType: "BOOL"},

		//
		{Name: "CVFaulted", DataType: "BOOL"},

		//
		{Name: "HandFBFaulted", DataType: "BOOL"},

		//
		{Name: "PVSpanInv", DataType: "BOOL"},

		//
		{Name: "SPProgInv", DataType: "BOOL"},

		//
		{Name: "SPOperInv", DataType: "BOOL"},

		//
		{Name: "SPCascadeInv", DataType: "BOOL"},

		//
		{Name: "SPLimitsInv", DataType: "BOOL"},

		//
		{Name: "RatioProgInv", DataType: "BOOL"},

		//
		{Name: "RatioOperInv", DataType: "BOOL"},

		//
		{Name: "RatioLimitsInv", DataType: "BOOL"},

		//
		{Name: "CVProgInv", DataType: "BOOL"},

		//
		{Name: "CVOperInv", DataType: "BOOL"},

		//
		{Name: "CVOverrideInv", DataType: "BOOL"},

		//
		{Name: "CVPreviousInv", DataType: "BOOL"},

		//
		{Name: "CVEUSpanInv", DataType: "BOOL"},

		//
		{Name: "CVLimitsInv", DataType: "BOOL"},

		//
		{Name: "CVROCLimitInv", DataType: "BOOL"},

		//
		{Name: "FFInv", DataType: "BOOL"},

		//
		{Name: "FFPreviousInv", DataType: "BOOL"},

		//
		{Name: "HandFBInv", DataType: "BOOL"},

		//
		{Name: "PGainInv", DataType: "BOOL"},

		//
		{Name: "IGainInv", DataType: "BOOL"},

		//
		{Name: "DGainInv", DataType: "BOOL"},

		//
		{Name: "ZCDeadbandInv", DataType: "BOOL"},

		//
		{Name: "PVDeadbandInv", DataType: "BOOL"},

		//
		{Name: "PVROCLimitsInv", DataType: "BOOL"},

		//
		{Name: "DevHLLimitsInv", DataType: "BOOL"},

		//
		{Name: "DevDeadbandInv", DataType: "BOOL"},

		//
		{Name: "AtuneDataInv", DataType: "BOOL"},

		//
		{Name: "TimingModeInv", DataType: "BOOL"},

		//
		{Name: "RTSMissed", DataType: "BOOL"},

		//
		{Name: "RTSTimeInv", DataType: "BOOL"},

		//
		{Name: "RTSTimeStampInv", DataType: "BOOL"},

		//
		{Name: "DeltaTInv", DataType: "BOOL"},
	},
}
