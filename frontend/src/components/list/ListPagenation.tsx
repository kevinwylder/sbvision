import * as React from 'react';

interface ListPagenationProps {
    maxTotal: number
    pageSize: number
    start: number
    end: number
    onPageSelected: (offset: number) => void
}

const disabled = "#9a9a9a";

export function ListPagenation(props: ListPagenationProps) {
    let pageText = (props.pageSize <= 1) ? props.start + 1 + "" : props.start + 1 + " to " + props.end;
    const canGoBack = () => props.start > 0;
    const canGoForward = () => props.end < props.maxTotal;
    return <div
        className="list-pagenation">
            <div
                style={{
                    color: canGoBack() ? "black" : disabled
                }}
                onClick={() => {
                    if (canGoBack())  {
                        props.onPageSelected(Math.max(0, props.start - props.pageSize));
                    }
                }}
            >Back</div>
            <div> {pageText} </div>
            <div
                style={{
                    color: canGoForward() ? "black" : disabled
                }}
                onClick={() => {
                    if (canGoForward()) {
                        props.onPageSelected(Math.min(props.start + props.pageSize, props.maxTotal - props.pageSize));
                    }
                }}>Forward</div>
    </div>
}
