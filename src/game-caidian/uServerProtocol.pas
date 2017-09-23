unit uServerProtocol;

interface
uses uTypeCommon;

const
    CM_AGENT_AUTH = 0;
    SM_AGENT_AUTH = 1;
    CM_AGENT_OPERATE = 2;
    SM_AGENT_OPERATE = 3;
    SM_GAMESTATUS = 4;
    SM_BROADCAST_DICE = 5;
    SM_AGENT_GIFT = 6;
    SM_AGENT_BECOMEBANKER = 7;
    SM_BROADCAST_WANGER = 8;
    SM_AGENT_M2BROADCAST = 9;
    SM_AGENT_PREDICERESULT = 10;
    SM_GAME_CONFIG = 11;


type
    TAgentAuthorizeReq = packed record
        Head: TMsgHead;
        LicKey: array[0..31] of byte;
    end;

    TAgentAuthorizeResp = packed record
        Head: TMsgHead;
        Err: byte;
        EncKey: cardinal;
    end;
    PAgentAuthorizeResp = ^TAgentAuthorizeResp;


//********************
    TAgentOperateReq = packed record
        Head: TMsgHead;
        Operate: byte;
        Reserved: uint64;
        OpGold: uint64;
        PlayName: array[0..30] of byte;
        NickName: array[0..30] of byte;
        Targ: byte;
    end;

    TAgentOperateResp = packed record
        Head: TMsgHead;
        Operate: byte;
        Reserved: uint64;
        PlayName: array[0..30] of byte;
        NickName: array[0..30] of byte;
        Err: byte;
    end;
    PAgentOperateResp = ^TAgentOperateResp;


//************************
    TGameStateResp = packed record
        Head: TMsgHead;
        Step: byte;
        Stay: byte;
    end;
    PGameStateResp = ^TGameStateResp;


//*********************
    TDiceResult = packed record
        DiceVal: array[0..2] of byte;
        Result: byte;
    end;

    TDiceResultResp = packed record
        Head: TMsgHead;
        Dice: TDiceResult;
    end;
    PDiceResultResp = ^TDiceResultResp;


//*************************

    TAgentGiftResp = packed record
        Head: TMsgHead;
        PlayName: array[0..30] of byte;
        NickName: array[0..30] of byte;
        Reserved: uint64;
        Gold: uint64;
    end;
    PAgentGiftResp = ^TAgentGiftResp;

//*************************

    TAgentBecomeBankerResp = packed record
        Head: TMsgHead;
        Region: array[0..30] of byte;
        PlayName: array[0..30] of byte;
        NickName: array[0..30] of byte;
        Gold: uint64;
        Reserved: uint64;
    end;
    PAgentBecomeBankerResp = ^TAgentBecomBankerResp;

//*************************

    TWangerNotifyResp = packed record
        Head: TMsgHead;
        BigGold: uint64;
        BigLim: uint64;
        SmallGold: uint64;
        SmallLim: uint64;
    end;
    PWangerNotifyResp = ^TWangerNotifyResp;


//*****************unimplemented yet************************

    TMsgContent = packed record
        Content: array[0..128] of byte;
        FColor: byte;
        BColor: byte;
        Position: byte;
    end;
    TBroadcastM2Req = packed record
        Head: TMsgHead;
        Content: TMsgContent;
    end;
    PBroadcastM2Req = ^TBroadcastM2Req;

//********************
implementation



end.