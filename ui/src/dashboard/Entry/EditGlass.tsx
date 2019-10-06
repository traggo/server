import * as React from 'react';
import EditIcon from '@material-ui/icons/Edit';
import Typography from '@material-ui/core/Typography';
import {Center} from '../../common/Center';
import DeleteIcon from '@material-ui/icons/Delete';
import {Button} from '@material-ui/core';

interface EditGlassProps {
    doEdit: (elm: HTMLElement) => void;
    doDelete: () => void;
}

export const EditGlass: React.FC<EditGlassProps> = ({doEdit, doDelete}) => {
    const [confirm, askConfirmation] = React.useState(false);
    return (
        <div
            style={{
                position: 'absolute',
                top: 0,
                bottom: 0,
                left: 0,
                right: 0,
                background: 'rgba(110,110,110,.5)',
            }}>
            {confirm ? (
                <div>
                    <Typography>Confirm Deletion</Typography>
                    <Button onClick={() => askConfirmation(false)}>Cancel</Button>
                    <Button color={'primary'} onClick={doDelete}>
                        Yeah, Delete
                    </Button>
                </div>
            ) : (
                <>
                    <div style={{height: '30%', display: 'flex', cursor: 'pointer'}}>
                        <div
                            style={{height: '100%', width: '65%', background: 'rgba(255,255,255,.4)'}}
                            onClick={(e) => doEdit(e.currentTarget)}>
                            <Center>
                                <EditIcon />
                            </Center>
                        </div>
                        <div
                            style={{height: '100%', width: '35%', background: 'rgba(192, 57, 43, .4)'}}
                            onClick={() => askConfirmation(true)}>
                            <Center>
                                <DeleteIcon fontSize={'small'} />
                            </Center>
                        </div>
                    </div>
                    <div style={{height: '70%', cursor: 'grab'}}>
                        <Center>
                            <div style={{background: 'rgba(120,120,120, .8)', padding: 10}}>
                                <Typography variant={'h5'} style={{textAlign: 'center', width: '100%'}}>
                                    drag here
                                </Typography>
                                <Typography variant={'body1'} style={{textAlign: 'center'}}>
                                    resize on corner
                                </Typography>
                            </div>
                        </Center>
                    </div>
                </>
            )}
        </div>
    );
};
