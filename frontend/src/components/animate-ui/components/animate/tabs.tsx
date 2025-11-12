import * as React from 'react';

import {
  Tabs as TabsPrimitive,
  TabsList as TabsListPrimitive,
  TabsTrigger as TabsTriggerPrimitive,
  TabsContent as TabsContentPrimitive,
  TabsContents as TabsContentsPrimitive,
  TabsHighlight as TabsHighlightPrimitive,
  TabsHighlightItem as TabsHighlightItemPrimitive,
  type TabsProps as TabsPrimitiveProps,
  type TabsListProps as TabsListPrimitiveProps,
  type TabsTriggerProps as TabsTriggerPrimitiveProps,
  type TabsContentProps as TabsContentPrimitiveProps,
  type TabsContentsProps as TabsContentsPrimitiveProps,
} from '@/components/animate-ui/primitives/animate/tabs';
import { cn } from '@/lib/utils';

type TabsProps = TabsPrimitiveProps;

function Tabs({ className, ...props }: TabsProps) {
  return (
    <TabsPrimitive
      className={cn('flex flex-col gap-0', className)}
      {...props}
    />
  );
}

type TabsListProps = TabsListPrimitiveProps;

function TabsList({ className, ...props }: TabsListProps) {
  return (
    <TabsListPrimitive
      className={cn(
        'inline-flex h-12 items-center justify-center rounded-lg p-1 mx-auto bg-[#ffb947] relative',
        className,
      )}
      {...props}
    />
  );
}

type TabsTriggerProps = TabsTriggerPrimitiveProps;

function TabsTrigger({ className, ...props }: TabsTriggerProps) {
  return (
    <TabsTriggerPrimitive
      className={cn(
        "data-[state=active]:bg-[#F0711A] data-[state=active]:text-white data-[state=active]:shadow-md focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:outline-ring inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md w-full px-4 py-2 text-base font-medium whitespace-nowrap transition-all duration-300 ease-in-out focus-visible:ring-[3px] focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 text-[#161B27]",
        className,
      )}
      {...props}
    />
  );
}

type TabsContentsProps = TabsContentsPrimitiveProps;

function TabsContents(props: TabsContentsProps) {
  return <TabsContentsPrimitive {...props} />;
}

type TabsContentProps = TabsContentPrimitiveProps;

function TabsContent({ className, ...props }: TabsContentProps) {
  return (
    <TabsContentPrimitive
      className={cn('outline-none', className)}
      {...props}
    />
  );
}

const TabsHighlight = TabsHighlightPrimitive;
const TabsHighlightItem = TabsHighlightItemPrimitive;

export {
  Tabs,
  TabsList,
  TabsTrigger,
  TabsContents,
  TabsContent,
  TabsHighlight,
  TabsHighlightItem,
  type TabsProps,
  type TabsListProps,
  type TabsTriggerProps,
  type TabsContentsProps,
  type TabsContentProps,
};
