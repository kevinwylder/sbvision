import { Quaternion } from './math';

export interface Label{
    input: ModelInput
    output: ModelOutput
}

export interface ModelInput{
    width: number
    height: number
    data: string
}

export interface ModelOutput{
    isSkateboard: number
    rotation: Quaternion
}
